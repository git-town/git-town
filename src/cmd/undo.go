package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
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
			return undo(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func undo(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	_, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  &repo,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	runState, err := runstate.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsUnfinished() {
		return fmt.Errorf(messages.UndoNothingToDo)
	}
	undoRunState := runState.CreateUndoRunState()
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &undoRunState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}
