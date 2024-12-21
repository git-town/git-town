package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/cmd/ship"
	"github.com/git-town/git-town/v17/internal/cmd/sync"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/slice"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/hosting"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/undo/undoconfig"
	"github.com/git-town/git-town/v17/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/program"
	"github.com/git-town/git-town/v17/internal/vm/runstate"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

const hoistDesc = "Extract a branch from a stack, making it a top-level branch"

const hoistHelp = `
Assume you have a stack. One of the branches in your stack makes changes that don't require the changes made by its ancestor branches. The "hoist" command removes this branch from the stack and makes it a stand-alone top-level branch that ships directly into your main branch. This allows you to get your changes reviewed and shipped concurrently rather than sequentially.`

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
	// branchToHoistInfo        gitdomain.BranchInfo
	// branchToHoistType        configdomain.BranchType
	// branchWhenDone           gitdomain.LocalBranchName
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	connector        Option[hostingdomain.Connector]
	dialogTestInputs components.TestInputs
	// dryRun                   configdomain.DryRun
	hasOpenChanges              bool
	initialBranch               gitdomain.LocalBranchName
	initialBranchContainsMerges bool
	initialBranchType           configdomain.BranchType
	nonExistingBranches         gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	// previousBranch           Option[gitdomain.LocalBranchName]
	// proposalsOfChildBranches []hostingdomain.Proposal
	stashSize gitdomain.StashSize
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
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
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
	branchWhenDone := determineBranchWhenDone(branchWhenDoneArgs{
		branchNameToHoist: branchNameToHoist,
		branches:          branchesSnapshot.Branches,
		initialBranch:     initialBranch,
		mainBranch:        validatedConfig.ValidatedConfigData.MainBranch,
		previousBranch:    previousBranchOpt,
	})
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
		branchWhenDone:           branchWhenDone,
		branchesSnapshot:         branchesSnapshot,
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
	switch data.branchToHoistType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		hoistFeatureBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch:
		hoistLocalBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToHoistType))
	}
	localBranchNameToHoist, hasLocalBranchToHoist := data.branchToHoistInfo.LocalName.Get()
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         hasLocalBranchToHoist && data.initialBranch != localBranchNameToHoist && data.hasOpenChanges,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch, Some(data.initialBranch)},
	})
	return prog.Immutable(), undoProg.Immutable()
}

func hoistFeatureBranch(prog, finalUndoProgram Mutable[program.Program], data hoistData) {
	trackingBranchToHoist, hasTrackingBranchToHoist := data.branchToHoistInfo.RemoteName.Get()
	if data.branchToHoistInfo.SyncStatus != gitdomain.SyncStatusHoistdAtRemote && hasTrackingBranchToHoist && data.config.NormalConfig.IsOnline() {
		ship.UpdateChildBranchProposalsToGrandParent(prog.Value, data.proposalsOfChildBranches)
		prog.Value.Add(&opcodes.BranchTrackingHoist{Branch: trackingBranchToHoist})
	}
	hoistLocalBranch(prog, finalUndoProgram, data)
}

func hoistLocalBranch(prog, finalUndoProgram Mutable[program.Program], data hoistData) {
	if localBranchToHoist, hasLocalBranchToHoist := data.branchToHoistInfo.LocalName.Get(); hasLocalBranchToHoist {
		if data.initialBranch == localBranchToHoist {
			if data.hasOpenChanges {
				prog.Value.Add(&opcodes.ChangesStage{})
				prog.Value.Add(&opcodes.CommitWithMessage{
					AuthorOverride: None[gitdomain.Author](),
					Message:        gitdomain.CommitMessage("Committing open changes on hoistd branch"),
				})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Value.Add(&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localBranchToHoist})
				finalUndoProgram.Value.Add(&opcodes.UndoLastCommit{})
			}
		}
		// hoist the commits of this branch from all descendents
		if data.config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase {
			descendents := data.config.NormalConfig.Lineage.Descendants(localBranchToHoist)
			for _, descendent := range descendents {
				if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
					sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
						Ancestor:          localBranchToHoist.BranchName(),
						Branch:            descendent,
						HasTrackingBranch: branchInfo.HasTrackingBranch(),
						Program:           prog,
						RebaseOnto:        data.config.ValidatedConfigData.MainBranch,
					})
				}
			}
		}
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.branchWhenDone})
		prog.Value.Add(&opcodes.BranchLocalHoist{
			Branch: localBranchToHoist,
		})
		if data.dryRun.IsFalse() {
			sync.RemoveBranchConfiguration(sync.RemoveBranchConfigurationArgs{
				Branch:  localBranchToHoist,
				Lineage: data.config.NormalConfig.Lineage,
				Program: prog,
			})
		}
	}
}

func validateHoistData(data hoistData) error {
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
