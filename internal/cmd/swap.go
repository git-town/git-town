package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/cmd/sync"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/execute"
	"github.com/git-town/git-town/v18/internal/forge"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/gohacks/slice"
	"github.com/git-town/git-town/v18/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/undo/undoconfig"
	"github.com/git-town/git-town/v18/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v18/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v18/internal/vm/opcodes"
	"github.com/git-town/git-town/v18/internal/vm/program"
	"github.com/git-town/git-town/v18/internal/vm/runstate"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

const swapDesc = "Swap this branch with the one ahead of it in the stack"

const swapHelp = `
The "swap" command moves this branch one position forward in the stack.`

const swapCommandName = "swap"

func swapCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     swapCommandName,
		Args:    cobra.NoArgs,
		Short:   swapDesc,
		GroupID: "stack",
		Long:    cmdhelpers.Long(swapDesc, swapHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeSwap(args, dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwap(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSwapData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateSwapData(data)
	if err != nil {
		return err
	}
	runProgram := detachProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               detachCommandName,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type detachData struct {
	branchToSwapContainsMerges bool
	branchToSwapInfo           gitdomain.BranchInfo
	branchToSwapName           gitdomain.LocalBranchName
	branchToSwapType           configdomain.BranchType
	branchesSnapshot           gitdomain.BranchesSnapshot
	children                   []detachChildBranch
	config                     config.ValidatedConfig
	connector                  Option[forgedomain.Connector]
	dialogTestInputs           components.TestInputs
	dryRun                     configdomain.DryRun
	hasOpenChanges             bool
	initialBranch              gitdomain.LocalBranchName
	nonExistingBranches        gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch               gitdomain.LocalBranchName
	previousBranch             Option[gitdomain.LocalBranchName]
	stashSize                  gitdomain.StashSize
}

type detachChildBranch struct {
	info     gitdomain.BranchInfo
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineDetachData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data detachData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchNameToDetach := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToDetachInfo, hasBranchToDetachInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToDetach).Get()
	if !hasBranchToDetachInfo {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToDetach)
	}
	if branchToDetachInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToDetach)
	}
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchTypeToDetach := validatedConfig.BranchType(branchNameToDetach)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToDetach).Get()
	if !hasParentBranch {
		return data, false, errors.New(messages.DetachNoParent)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchNameToDetach)
	children := make([]detachChildBranch, len(childBranches))
	for c, childBranch := range childBranches {
		proposal := None[forgedomain.Proposal]()
		if connector, hasConnector := connector.Get(); hasConnector {
			if findProposal, canFindProposal := connector.FindProposalFn().Get(); canFindProposal {
				proposal, err = findProposal(childBranch, initialBranch)
				if err != nil {
					return data, false, err
				}
			}
		}
		childInfo, has := branchesSnapshot.Branches.FindByLocalName(childBranch).Get()
		if !has {
			return data, false, fmt.Errorf("cannot find branch info for %q", childBranch)
		}
		children[c] = detachChildBranch{
			info:     *childInfo,
			name:     childBranch,
			proposal: proposal,
		}
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return detachData{
		branchToDetachContainsMerges: false, // TODO: determine the actual data
		branchToDetachInfo:           *branchToDetachInfo,
		branchToDetachName:           branchNameToDetach,
		branchToDetachType:           branchTypeToDetach,
		branchesSnapshot:             branchesSnapshot,
		children:                     children,
		config:                       validatedConfig,
		connector:                    connector,
		dialogTestInputs:             dialogTestInputs,
		dryRun:                       dryRun,
		hasOpenChanges:               repoStatus.OpenChanges,
		initialBranch:                initialBranch,
		nonExistingBranches:          nonExistingBranches,
		parentBranch:                 parentBranch,
		previousBranch:               previousBranchOpt,
		stashSize:                    stashSize,
	}, false, nil
}

func detachProgram(data detachData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages)
	prog.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseAgainst: data.parentBranch.BranchName(),
			BranchToRebaseOnto:    data.config.ValidatedConfigData.MainBranch,
			Upstream:              None[gitdomain.LocalBranchName](),
		},
	)
	if data.branchToDetachInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true},
		)
	}
	lastParent := data.parentBranch
	descendents := data.config.NormalConfig.Lineage.Descendants(data.branchToDetachName)
	for _, descendent := range descendents {
		if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
			sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
				Ancestor:          data.branchToDetachName.BranchName(),
				Branch:            descendent,
				HasTrackingBranch: branchInfo.HasTrackingBranch(),
				Program:           prog,
				RebaseOnto:        lastParent,
			})
			if branchInfo.HasTrackingBranch() {
				prog.Value.Add(
					&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true},
				)
			}
			lastParent = descendent
		}
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch})
	if data.dryRun.IsFalse() {
		prog.Value.Add(
			&opcodes.LineageParentSet{
				Branch: data.branchToDetachName,
				Parent: data.config.ValidatedConfigData.MainBranch,
			},
		)
		for _, child := range data.config.NormalConfig.Lineage.Children(data.branchToDetachName) {
			prog.Value.Add(
				&opcodes.LineageParentSet{
					Branch: child,
					Parent: data.parentBranch,
				},
			)
		}
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable()
}

func validateDetachData(data detachData) error {
	switch data.branchToDetachInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusBehind:
		return errors.New(messages.DetachNeedsSync)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.DetachOtherWorkTree, data.branchToDetachName)
	case gitdomain.SyncStatusRemoteOnly:
		return errors.New(messages.DetachRemoteBranch)
	}
	if data.branchToDetachContainsMerges {
		return fmt.Errorf(messages.BranchContainsMergeCommits, data.initialBranch)
	}
	switch data.branchToDetachType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.DetachUnsupportedBranchType, data.branchToDetachType)
	}
	for _, child := range data.children {
		switch child.info.SyncStatus {
		case
			gitdomain.SyncStatusAhead,
			gitdomain.SyncStatusLocalOnly,
			gitdomain.SyncStatusUpToDate:
		case
			gitdomain.SyncStatusBehind,
			gitdomain.SyncStatusDeletedAtRemote,
			gitdomain.SyncStatusNotInSync,
			gitdomain.SyncStatusRemoteOnly:
			return errors.New(messages.DetachNeedsSync)
		case gitdomain.SyncStatusOtherWorktree:
			return fmt.Errorf(messages.DetachOtherWorkTree, child.name)
		}
	}
	return nil
}
