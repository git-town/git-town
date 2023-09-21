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
	"github.com/git-town/git-town/v9/src/undo"
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
			return runUndo(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runUndo(debug bool) error {
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
	config, initialBranchesSnapshot, initialStashSnaphot, lineage, err := determineUndoConfig(&repo)
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
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnaphot,
	})
}

type undoConfig struct {
	hasOpenChanges          bool
	mainBranch              domain.LocalBranchName
	initialBranchesSnapshot undo.BranchesSnapshot
	previousBranch          domain.LocalBranchName
}

func determineUndoConfig(repo *execute.OpenRepoResult) (*undoConfig, undo.BranchesSnapshot, undo.StashSnapshot, config.Lineage, error) {
	lineage := repo.Runner.Config.Lineage()
	_, initialBranchesSnapshot, initialStashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, initialBranchesSnapshot, initialStashSnapshot, lineage, err
	}
	mainBranch := repo.Runner.Backend.Config.MainBranch()
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
	if err != nil {
		return nil, initialBranchesSnapshot, initialStashSnapshot, lineage, err
	}
	return &undoConfig{
		hasOpenChanges:          hasOpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		mainBranch:              mainBranch,
		previousBranch:          previousBranch,
	}, initialBranchesSnapshot, initialStashSnapshot, lineage, nil
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
	fmt.Println("1111111111111111111111")
	fmt.Printf("%#v\n", undoRunState)
	err = undoRunState.RunStepList.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranchesSnapshot.Active,
		PreviousBranch:   config.previousBranch,
	})
	return undoRunState, err
}
