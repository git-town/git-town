package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineKillConfig(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateKillConfig(config)
	if err != nil {
		return err
	}
	steps, finalUndoProgram := killProgram(config)
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
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type killConfig struct {
	*configdomain.FullConfig
	branchNameToKill gitdomain.BranchInfo
	branchTypeToKill configdomain.BranchType
	branchWhenDone   gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	parentBranch     Option[gitdomain.LocalBranchName]
	previousBranch   gitdomain.LocalBranchName
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (*killConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
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
	if branchToKill.IsLocal() {
		err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
			BranchesToVerify: gitdomain.LocalBranchNames{branchToKill.LocalName},
			Config:           repo.Runner.Config,
			DefaultChoice:    repo.Runner.Config.FullConfig.MainBranch,
			DialogTestInputs: &dialogTestInputs,
			LocalBranches:    branchesSnapshot.Branches,
			MainBranch:       repo.Runner.Config.FullConfig.MainBranch,
			Runner:           repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
	}
	branchTypeToKill := repo.Runner.Config.FullConfig.BranchType(branchNameToKill)
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	var branchWhenDone gitdomain.LocalBranchName
	if branchNameToKill == branchesSnapshot.Active {
		if previousBranch == branchesSnapshot.Active {
			branchWhenDone = repo.Runner.Config.FullConfig.MainBranch
		} else {
			branchWhenDone = previousBranch
		}
	} else {
		branchWhenDone = branchesSnapshot.Active
	}
	parentBranch := repo.Runner.Config.FullConfig.Lineage.Parent(branchToKill.LocalName)
	return &killConfig{
		FullConfig:       &repo.Runner.Config.FullConfig,
		branchNameToKill: branchToKill,
		branchTypeToKill: branchTypeToKill,
		branchWhenDone:   branchWhenDone,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		parentBranch:     parentBranch,
		previousBranch:   previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func killProgram(config *killConfig) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	switch config.branchTypeToKill {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch:
		killFeatureBranch(&prog, &finalUndoProgram, config)
	case configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch:
		killLocalBranch(&prog, &finalUndoProgram, config)
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", config.branchTypeToKill))
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.initialBranch != config.branchNameToKill.LocalName && config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.initialBranch},
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, config *killConfig) {
	if config.branchNameToKill.HasTrackingBranch() && config.IsOnline() {
		prog.Add(&opcodes.DeleteTrackingBranch{Branch: config.branchNameToKill.RemoteName})
	}
	killLocalBranch(prog, finalUndoProgram, config)
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killLocalBranch(prog *program.Program, finalUndoProgram *program.Program, config *killConfig) {
	if config.initialBranch == config.branchNameToKill.LocalName {
		if config.hasOpenChanges {
			prog.Add(&opcodes.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			prog.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoProgram.Add(&opcodes.Checkout{Branch: config.branchNameToKill.LocalName})
			finalUndoProgram.Add(&opcodes.UndoLastCommit{})
		}
		prog.Add(&opcodes.Checkout{Branch: config.branchWhenDone})
	}
	prog.Add(&opcodes.DeleteLocalBranch{Branch: config.branchNameToKill.LocalName})
	if parentBranch, hasParentBranch := config.parentBranch.Get(); hasParentBranch && !config.dryRun {
		sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
			Branch:  config.branchNameToKill.LocalName,
			Lineage: config.Lineage,
			Parent:  parentBranch,
			Program: prog,
		})
	}
}

func validateKillConfig(killConfig *killConfig) error {
	switch killConfig.branchTypeToKill {
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.KillCannotKillMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.KillCannotKillPerennialBranches)
	}
	panic(fmt.Sprintf("unhandled branch type: %s", killConfig.branchTypeToKill))
}
