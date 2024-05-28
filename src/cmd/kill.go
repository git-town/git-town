package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
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

func executeKill(args []string, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineKillData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateKillData(data)
	if err != nil {
		return err
	}
	steps, finalUndoProgram := killProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "kill",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            steps,
		FinalUndoProgram:      finalUndoProgram,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type killData struct {
	branchToKillInfo gitdomain.BranchInfo
	branchToKillType configdomain.BranchType
	branchWhenDone   gitdomain.LocalBranchName
	config           config.ValidatedConfig
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	parentBranch     Option[gitdomain.LocalBranchName]
	previousBranch   gitdomain.LocalBranchName
}

func determineKillData(args []string, repo execute.OpenRepoResult, dryRun, verbose bool) (*killData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToKill, hasBranchToKill := branchesSnapshot.Branches.FindByLocalName(branchNameToKill).Get()
	if !hasBranchToKill {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, branchesSnapshot, stashSize, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesToKill := gitdomain.LocalBranchNames{branchNameToKill}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchesToKill,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	branchTypeToKill := validatedConfig.Config.BranchType(branchNameToKill)
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return nil, branchesSnapshot, stashSize, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	var branchWhenDone gitdomain.LocalBranchName
	if branchNameToKill == initialBranch {
		if previousBranch == initialBranch {
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
	return &killData{
		branchToKillInfo: branchToKill,
		branchToKillType: branchTypeToKill,
		branchWhenDone:   branchWhenDone,
		config:           validatedConfig,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		parentBranch:     parentBranch,
		previousBranch:   previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func killProgram(data *killData) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	switch data.branchToKillType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch:
		killFeatureBranch(&prog, &finalUndoProgram, data)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch:
		killLocalBranch(&prog, &finalUndoProgram, data)
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToKillType))
	}
	localBranchNameToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get()
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         hasLocalBranchToKill && data.initialBranch != localBranchNameToKill && data.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{data.previousBranch, data.initialBranch},
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, data *killData) {
	trackingBranchToKill, hasTrackingBranchToKill := data.branchToKillInfo.RemoteName.Get()
	if data.branchToKillInfo.SyncStatus != gitdomain.SyncStatusDeletedAtRemote && hasTrackingBranchToKill && data.config.Config.IsOnline() {
		prog.Add(&opcodes.DeleteTrackingBranch{Branch: trackingBranchToKill})
	}
	killLocalBranch(prog, finalUndoProgram, data)
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killLocalBranch(prog, finalUndoProgram *program.Program, data *killData) {
	if localBranchToKill, hasLocalBranchToKill := data.branchToKillInfo.LocalName.Get(); hasLocalBranchToKill {
		if data.initialBranch == localBranchToKill {
			if data.hasOpenChanges {
				prog.Add(&opcodes.CommitOpenChanges{})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Add(&opcodes.Checkout{Branch: localBranchToKill})
				finalUndoProgram.Add(&opcodes.UndoLastCommit{})
			}
			prog.Add(&opcodes.Checkout{Branch: data.branchWhenDone})
		}
		prog.Add(&opcodes.DeleteLocalBranch{Branch: localBranchToKill})
		if parentBranch, hasParentBranch := data.parentBranch.Get(); hasParentBranch && !data.dryRun {
			sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
				Branch:  localBranchToKill,
				Lineage: data.config.Config.Lineage,
				Parent:  parentBranch,
				Program: prog,
			})
		}
	}
}

func validateKillData(data *killData) error {
	switch data.branchToKillType {
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.KillCannotKillMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.KillCannotKillPerennialBranches)
	}
	panic(fmt.Sprintf("unhandled branch type: %s", data.branchToKillType))
}
