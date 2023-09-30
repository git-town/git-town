package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/spf13/cobra"
)

const undoDesc = "Undoes the last run git-town command"

func undoCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    long(undoDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUndo(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executeUndo(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialStashSnaphot, lineage, err := determineUndoConfig(&repo)
	if err != nil {
		return err
	}
	undoRunState, err := determineUndoRunState(config, &repo)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &undoRunState,
		Run:                     &repo.Runner,
		Connector:               nil,
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

func determineUndoConfig(repo *execute.OpenRepoResult) (*undoConfig, domain.StashSnapshot, config.Lineage, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), lineage, err
	}
	_, initialBranchesSnapshot, initialStashSnapshot, _, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
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
	hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	return &undoConfig{
		hasOpenChanges:          hasOpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		mainBranch:              mainBranch,
		previousBranch:          previousBranch,
		pushHook:                pushHook,
	}, initialStashSnapshot, lineage, nil
}

func determineUndoRunState(config *undoConfig, repo *execute.OpenRepoResult) (runstate.RunState, error) {
	runState, err := persistence.Load(repo.RootDir)
	if err != nil {
		return runstate.RunState{}, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsUnfinished() {
		return runstate.RunState{}, fmt.Errorf(messages.UndoNothingToDo)
	}
	undoRunState := runState.CreateUndoRunState()
	err = undoRunState.RunStepList.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranchesSnapshot.Active,
		PreviousBranch:   config.previousBranch,
	})
	if err != nil {
		return runstate.RunState{}, err
	}
	// If the command to undo failed and was continued,
	// there might be steps in the undo stack that became obsolete
	// when the command was continued.
	// Example: the command stashed away uncommitted changes,
	// failed, and remembered in the undo list to pop the stack.
	// When continuing, it finishes and pops the stack as part of the continue list.
	// When we run undo now, it still wants to pop the stack even though that was already done.
	// This seems to apply only to popping the stack and switching back to the initial branch.
	// Hence we consolidate this step type here.
	undoRunState.RunStepList.List = slice.LowerAll[steps.Step](undoRunState.RunStepList.List, &steps.RestoreOpenChangesStep{})
	undoRunState.RunStepList.RemoveAllButLast("*steps.CheckoutIfExistsStep")
	return undoRunState, err
}
