package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/runstate"
	"github.com/git-town/git-town/v10/src/vm/statefile"
	"github.com/spf13/cobra"
)

const undoDesc = "Undoes the last run git-town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    long(undoDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUndo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialStashSnaphot, lineage, err := determineUndoConfig(repo, verbose)
	if err != nil {
		return err
	}
	undoRunState, err := determineUndoRunState(config, repo)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &undoRunState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		Lineage:                 lineage,
		NoPushHook:              !config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: config.initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnaphot,
	})
}

type undoConfig struct {
	hasOpenChanges          bool
	mainBranch              domain.LocalBranchName
	initialBranchesSnapshot domain.BranchesSnapshot
	previousBranch          domain.LocalBranchName
	pushHook                bool
}

func determineUndoConfig(repo *execute.OpenRepoResult, verbose bool) (*undoConfig, domain.StashSnapshot, config.Lineage, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), lineage, err
	}
	_, initialBranchesSnapshot, initialStashSnapshot, _, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	mainBranch := repo.Runner.Backend.Config.MainBranch()
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	return &undoConfig{
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		mainBranch:              mainBranch,
		previousBranch:          previousBranch,
		pushHook:                pushHook,
	}, initialStashSnapshot, lineage, nil
}

func determineUndoRunState(config *undoConfig, repo *execute.OpenRepoResult) (runstate.RunState, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.RunState{}, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsUnfinished() {
		return runstate.RunState{}, fmt.Errorf(messages.UndoNothingToDo)
	}
	undoRunState := runState.CreateUndoRunState()
	wrap(&undoRunState.RunProgram, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranchesSnapshot.Active,
		PreviousBranch:   config.previousBranch,
	})
	// If the command to undo failed and was continued,
	// there might be opcodes in the undo stack that became obsolete
	// when the command was continued.
	// Example: the command stashed away uncommitted changes,
	// failed, and remembered in the undo list to pop the stack.
	// When continuing, it finishes and pops the stack as part of the continue list.
	// When we run undo now, it still wants to pop the stack even though that was already done.
	// This seems to apply only to popping the stack and switching back to the initial branch.
	// Hence we consolidate these opcode types here.
	undoRunState.RunProgram.MoveToEnd(&opcode.RestoreOpenChanges{})
	undoRunState.RunProgram.RemoveAllButLast("*opcode.CheckoutIfExists")
	return undoRunState, err
}
