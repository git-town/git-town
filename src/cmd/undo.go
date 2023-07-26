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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	rootDir, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	_, _, err = execute.LoadBranches(&run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	runState, err := runstate.Load(rootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil || runState.IsUnfinished() {
		return fmt.Errorf(messages.UndoNothingToDo)
	}
	undoRunState := runState.CreateUndoRunState()
	return runstate.Execute(&undoRunState, &run, nil, rootDir)
}
