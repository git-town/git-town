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

const hoistDesc = "Extract a branch from a stack, making it a top-level branch"

const hoistHelp = `
Assume one of the branches in a stack makes changes that don't require the changes made by branches. This branch could  The "hoist" command removes this branch from the stack and makes it a stand-alone top-level branch that ships directly into your main branch. This allows you to get your changes reviewed and shipped concurrently rather than sequentially.`

const hoistCommandName = "hoist"

func hoistCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   hoistCommandName,
		Args:  cobra.NoArgs,
		Short: hoistDesc,
		Long:  cmdhelpers.Long(hoistDesc, hoistHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeHoist(args, dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHoist(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineHoistData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateHoistData(data)
	if err != nil {
		return err
	}
	runProgram := hoistProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               hoistCommandName,
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

type hoistData struct {
	branchToHoistContainsMerges bool
	branchToHoistInfo           gitdomain.BranchInfo
	branchToHoistName           gitdomain.LocalBranchName
	branchToHoistType           configdomain.BranchType
	branchesSnapshot            gitdomain.BranchesSnapshot
	children                    []hoistChildBranch
	config                      config.ValidatedConfig
	connector                   Option[forgedomain.Connector]
	dialogTestInputs            components.TestInputs
	dryRun                      configdomain.DryRun
	hasOpenChanges              bool
	initialBranch               gitdomain.LocalBranchName
	nonExistingBranches         gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch                gitdomain.LocalBranchName
	previousBranch              Option[gitdomain.LocalBranchName]
	stashSize                   gitdomain.StashSize
}

type hoistChildBranch struct {
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineHoistData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data hoistData, exit bool, err error) {
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
	branchNameToHoist := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToHoistInfo, hasBranchToHoistInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToHoist).Get()
	if !hasBranchToHoistInfo {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToHoist)
	}
	if branchToHoistInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToHoist)
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
	branchTypeToHoist := validatedConfig.BranchType(branchNameToHoist)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToHoist).Get()
	if !hasParentBranch {
		return data, false, errors.New("cannot hoist branch without parent")
	}
	childBranches := data.config.NormalConfig.Lineage.Children(initialBranch)
	children := make([]hoistChildBranch, len(childBranches))
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
		children[c] = hoistChildBranch{
			name:     childBranch,
			proposal: proposal,
		}
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return hoistData{
		branchToHoistContainsMerges: false, // TODO: determine the actual data
		branchToHoistInfo:           *branchToHoistInfo,
		branchToHoistName:           branchNameToHoist,
		branchToHoistType:           branchTypeToHoist,
		branchesSnapshot:            branchesSnapshot,
		children:                    children,
		config:                      validatedConfig,
		connector:                   connector,
		dialogTestInputs:            dialogTestInputs,
		dryRun:                      dryRun,
		hasOpenChanges:              repoStatus.OpenChanges,
		initialBranch:               initialBranch,
		nonExistingBranches:         nonExistingBranches,
		parentBranch:                parentBranch,
		previousBranch:              previousBranchOpt,
		stashSize:                   stashSize,
	}, false, nil
}

func hoistProgram(data hoistData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages)
	prog.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseAgainst: data.parentBranch.BranchName(),
			BranchToRebaseOnto:    data.config.ValidatedConfigData.MainBranch,
			Upstream:              None[gitdomain.LocalBranchName](),
		},
	)
	if data.branchToHoistInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true},
		)
	}
	lastParent := data.parentBranch
	descendents := data.config.NormalConfig.Lineage.Descendants(data.branchToHoistName)
	for _, descendent := range descendents {
		if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
			sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
				Ancestor:          data.branchToHoistName.BranchName(),
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
		data.config.NormalConfig.SetParent(data.branchToHoistName, data.config.ValidatedConfigData.MainBranch)
		for _, child := range data.config.NormalConfig.Lineage.Children(data.branchToHoistName) {
			data.config.NormalConfig.SetParent(child, data.parentBranch)
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

func validateHoistData(data hoistData) error {
	// TODO: ensure all branches are in sync or local only
	switch data.branchToHoistInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusBehind, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync:
		return fmt.Errorf("please sync your branches before hoisting")
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf("this branch cannot be hoisted because it is checked out in another worktree")
	case gitdomain.SyncStatusRemoteOnly:
		return fmt.Errorf("cannot hoist a remote branch")
	}
	if data.branchToHoistContainsMerges {
		return fmt.Errorf(messages.BranchContainsMergeCommits, data.initialBranch)
	}
	switch data.branchToHoistType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.HoistUnsupportedBranchType, data.branchToHoistType)
	}
	return nil
}
