package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/cmd/ship"
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
	runProgram, finalUndoProgram := hoistProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               hoistCommandName,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      finalUndoProgram,
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
		InitialBranch:           data.initialBranch.name,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type hoistData struct {
	initialBranch struct {
		containsMerges bool
		name           gitdomain.LocalBranchName
		info           gitdomain.BranchInfo
		branchType     configdomain.BranchType
	}
	branchesSnapshot gitdomain.BranchesSnapshot
	children         []struct {
		name     gitdomain.LocalBranchName
		proposal Option[forgedomain.Proposal]
	}
	config           config.ValidatedConfig
	connector        Option[forgedomain.Connector]
	dialogTestInputs components.TestInputs
	dryRun           configdomain.DryRun
	hasOpenChanges   bool
	// nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch struct {
		name     gitdomain.LocalBranchName
		proposal Option[forgedomain.Proposal]
	}
	previousBranch Option[gitdomain.LocalBranchName]
	stashSize      gitdomain.StashSize
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
	branchToHoist, hasBranchToHoist := branchesSnapshot.Branches.FindByLocalName(branchNameToHoist).Get()
	if !hasBranchToHoist {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToHoist)
	}
	if branchToHoist.SyncStatus == gitdomain.SyncStatusOtherWorktree {
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
	childBranches := data.config.NormalConfig.Lineage.Children(initialBranch)
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    repo.IsOffline,
		OldBranch:                  branchNameToHoist,
		OldBranchHasTrackingBranch: branchToHoist.HasTrackingBranch(),
	})
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return hoistData{
		branchToHoistInfo:        *branchToHoist,
		branchToHoistType:        branchTypeToHoist,
		branchesSnapshot:         branchesSnapshot,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connector,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		nonExistingBranches:      nonExistingBranches,
		previousBranch:           previousBranchOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, nil
}

func hoistProgram(data hoistData, finalMessages stringslice.Collector) (runProgram, finalUndoProgram program.Program) {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages)
	undoProg := NewMutable(&program.Program{})
	if isOmni, branchName, _ := data.branchToHoistInfo.IsOmniBranch(); isOmni {
		hoistFeatureBranch(prog, branchName, undoProg, data)
	} else if isLocalOnly, branchName := data.branchToHoistInfo.IsLocalOnlyBranch(); isLocalOnly {
		hoistLocalBranch(prog, branchName, undoProg, data)
	} else {
		// cannot hoist this branch
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable(), undoProg.Immutable()
}

func hoistFeatureBranch(prog Mutable[program.Program], branchName gitdomain.LocalBranchName, finalUndoProgram Mutable[program.Program], data hoistData) {
	// trackingBranchToHoist, hasTrackingBranchToHoist := data.branchToHoistInfo.RemoteName.Get()
	// if data.branchToHoistInfo.SyncStatus != gitdomain.SyncStatusHoistdAtRemote && hasTrackingBranchToHoist && data.config.NormalConfig.IsOnline() {
	// 	ship.UpdateChildBranchProposalsToGrandParent(prog.Value, data.proposalsOfChildBranches)
	// 	prog.Value.Add(&opcodes.BranchTrackingHoist{Branch: trackingBranchToHoist})
	// }
	// hoistLocalBranch(prog, finalUndoProgram, data)
}

func hoistLocalBranch(prog Mutable[program.Program], branchName gitdomain.LocalBranchName, finalUndoProgram Mutable[program.Program], data hoistData) {
	// make this branch a child of the main branch
	prog.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseAgainst: data.parentBranch.BranchName(),
			BranchToRebaseOnto:    data.config.ValidatedConfigData.MainBranch,
			Upstream:              None[gitdomain.LocalBranchName](),
		},
	)
	// hoist the commits of this branch from all descendents
	lastParent := data.parentBranch
	descendents := data.config.NormalConfig.Lineage.Descendants(branchName)
	for _, descendent := range descendents {
		if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
			sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
				Ancestor:          branchName.BranchName(),
				Branch:            descendent,
				HasTrackingBranch: branchInfo.HasTrackingBranch(),
				Program:           prog,
				RebaseOnto:        lastParent,
			})
			lastParent = descendent
		}
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch})
	if data.dryRun.IsFalse() {
		// update lineage
		data.config.NormalConfig.SetParent(branchName, data.config.ValidatedConfigData.MainBranch)
		for _, child := range data.config.NormalConfig.Lineage.Children(branchName) {
			data.config.NormalConfig.SetParent(child, branchName)
		}
	}
}

func validateHoistData(data hoistData) error {
	// TODO: ensure all branches are in sync or local only
	if data.initialBranchContainsMerges {
		return fmt.Errorf(messages.BranchContainsMergeCommits, data.initialBranch)
	}
	switch data.initialBranchType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.HoistUnsupportedBranchType, data.initialBranchType)
	}
	return nil
}
