package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/slice"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/sync"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	"github.com/spf13/cobra"
)

const killDesc = "Remove an obsolete feature branch"

const killHelp = `
Deletes the current or provided branch from the local and origin repositories. Does not delete perennial branches nor the main branch.`

func killCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: killDesc,
		Long:  cmdhelpers.Long(killDesc, killHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeKill(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeKill(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineKillData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateKillData(data)
	if err != nil {
		return err
	}
	runProgram, finalUndoProgram := killProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "kill",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      finalUndoProgram,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[hostingdomain.Connector](),
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

type killData struct {
	branchToKillInfo gitdomain.BranchInfo
	branchToKillType configdomain.BranchType
	branchWhenDone   gitdomain.LocalBranchName
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	dialogTestInputs components.TestInputs
	dryRun           configdomain.DryRun
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	parentBranch     Option[gitdomain.LocalBranchName]
	previousBranch   gitdomain.LocalBranchName
	stashSize        gitdomain.StashSize
}

func determineKillData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data killData, exit bool, err error) {
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
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToKill, hasBranchToKill := branchesSnapshot.Branches.FindByLocalName(branchNameToKill).Get()
	if !hasBranchToKill {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesToKill := gitdomain.LocalBranchNames{branchNameToKill}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchesToKill,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchTypeToKill := validatedConfig.Config.BranchType(branchNameToKill)
	previousBranch, hasPreviousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend).Get()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	var branchWhenDone gitdomain.LocalBranchName
	if branchNameToKill == initialBranch {
		if !hasPreviousBranch || previousBranch == initialBranch {
			branchWhenDone = validatedConfig.Config.MainBranch
		} else {
			branchWhenDone = previousBranch
		}
	} else {
		branchWhenDone = initialBranch
	}
	localBranchToKill, hasLocalBranchToKill := branchToKill.LocalName.Get()
	var parentBranch Option[gitdomain.LocalBranchName]
	if hasLocalBranchToKill {
		parentBranch = validatedConfig.Config.Lineage.Parent(localBranchToKill)
	} else {
		parentBranch = None[gitdomain.LocalBranchName]()
	}
	return killData{
		branchToKillInfo: *branchToKill,
		branchToKillType: branchTypeToKill,
		branchWhenDone:   branchWhenDone,
		branchesSnapshot: branchesSnapshot,
		config:           validatedConfig,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		parentBranch:     parentBranch,
		previousBranch:   previousBranch,
		stashSize:        stashSize,
	}, false, nil
}

func killProgram(data killData) (runProgram, finalUndoProgram program.Program) {
	prog := NewMutable(&program.Program{})
	undoProg := NewMutable(&program.Program{})
	switch data.branchToKillType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		killFeatureBranch(prog, undoProg, data)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch:
		killLocalBranch(prog, undoProg, data)
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToKillType))
	}
	localBranchNameToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get()
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         hasLocalBranchToKill && data.initialBranch != localBranchNameToKill && data.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{data.previousBranch, data.initialBranch},
	})
	return prog.Get(), undoProg.Get()
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog, finalUndoProgram Mutable[program.Program], data killData) {
	trackingBranchToKill, hasTrackingBranchToKill := data.branchToKillInfo.RemoteName.Get()
	if data.branchToKillInfo.SyncStatus != gitdomain.SyncStatusDeletedAtRemote && hasTrackingBranchToKill && data.config.Config.IsOnline() {
		prog.Value.Add(&opcodes.DeleteTrackingBranch{Branch: trackingBranchToKill})
	}
	killLocalBranch(prog, finalUndoProgram, data)
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killLocalBranch(prog, finalUndoProgram Mutable[program.Program], data killData) {
	if localBranchToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get(); hasLocalBranchToKill {
		if data.initialBranch == localBranchToKill {
			if data.hasOpenChanges {
				prog.Value.Add(&opcodes.CommitOpenChanges{
					AddAll:  true,
					Message: "Committing WIP for git town undo",
				})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Value.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Value.Add(&opcodes.Checkout{Branch: localBranchToKill})
				finalUndoProgram.Value.Add(&opcodes.UndoLastCommit{})
			}
			prog.Value.Add(&opcodes.Checkout{Branch: data.branchWhenDone})
		}
		prog.Value.Add(&opcodes.DeleteLocalBranch{Branch: localBranchToKill})
		if parentBranch, hasParentBranch := data.parentBranch.Get(); hasParentBranch && data.dryRun.IsFalse() {
			sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
				Branch:  localBranchToKill,
				Lineage: data.config.Config.Lineage,
				Parent:  parentBranch,
				Program: prog,
			})
		}
	}
}

func validateKillData(data killData) error {
	switch data.branchToKillType {
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.KillCannotKillMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.KillCannotKillPerennialBranches)
	}
	panic(fmt.Sprintf("unhandled branch type: %s", data.branchToKillType))
}
