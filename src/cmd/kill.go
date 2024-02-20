package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/sync"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/spf13/cobra"
)

const killDesc = "Removes an obsolete feature branch"

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
	steps, finalUndoProgram := killProgram(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		AfterBranchesSnapshot:  gitdomain.EmptyBranchesSnapshot(),
		AfterConfigSnapshot:    undoconfig.EmptyConfigSnapshot(),
		AfterStashSize:         0,
		BeforeBranchesSnapshot: initialBranchesSnapshot,
		BeforeConfigSnapshot:   repo.ConfigSnapshot,
		BeforeStashSize:        initialStashSize,
		Command:                "kill",
		DryRun:                 dryRun,
		RunProgram:             steps,
		FinalUndoProgram:       finalUndoProgram,
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
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

type killConfig struct {
	*configdomain.FullConfig
	branchToKill     gitdomain.BranchInfo
	branchWhenDone   gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (*killConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, stashSize, repoStatus, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FullConfig:            &repo.Runner.FullConfig,
		HandleUnfinishedState: false,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToKill := branchesSnapshot.Branches.FindByLocalName(branchNameToKill)
	if branchToKill == nil {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, branchesSnapshot, stashSize, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	if branchToKill.IsLocal() {
		err = execute.EnsureKnownBranchAncestry(branchToKill.LocalName, execute.EnsureKnownBranchAncestryArgs{
			Config:           &repo.Runner.FullConfig,
			AllBranches:      branchesSnapshot.Branches,
			DefaultBranch:    repo.Runner.MainBranch,
			DialogTestInputs: &dialogTestInputs,
			Runner:           repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
	}
	if !repo.Runner.IsFeatureBranch(branchToKill.LocalName) {
		return nil, branchesSnapshot, stashSize, false, errors.New(messages.KillOnlyFeatureBranches)
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	var branchWhenDone gitdomain.LocalBranchName
	if branchNameToKill == branchesSnapshot.Active {
		branchWhenDone = previousBranch
	} else {
		branchWhenDone = branchesSnapshot.Active
	}
	return &killConfig{
		FullConfig:       &repo.Runner.FullConfig,
		branchToKill:     *branchToKill,
		branchWhenDone:   branchWhenDone,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		previousBranch:   previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func (self killConfig) branchToKillParent() gitdomain.LocalBranchName {
	return self.Lineage.Parent(self.branchToKill.LocalName)
}

func killProgram(config *killConfig) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	killFeatureBranch(&prog, &finalUndoProgram, *config)
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.initialBranch != config.branchToKill.LocalName && config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.initialBranch},
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, config killConfig) {
	if config.branchToKill.HasTrackingBranch() && config.IsOnline() {
		prog.Add(&opcodes.DeleteTrackingBranch{Branch: config.branchToKill.RemoteName})
	}
	if config.initialBranch == config.branchToKill.LocalName {
		if config.hasOpenChanges {
			prog.Add(&opcodes.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			prog.Add(&opcodes.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoProgram.Add(&opcodes.Checkout{Branch: config.branchToKill.LocalName})
			finalUndoProgram.Add(&opcodes.UndoLastCommit{})
		}
		prog.Add(&opcodes.Checkout{Branch: config.branchWhenDone})
	}
	prog.Add(&opcodes.DeleteLocalBranch{Branch: config.branchToKill.LocalName, Force: false})
	if !config.dryRun {
		sync.RemoveBranchFromLineage(sync.RemoveBranchFromLineageArgs{
			Branch:  config.branchToKill.LocalName,
			Lineage: config.Lineage,
			Parent:  config.branchToKillParent(),
			Program: prog,
		})
	}
}
