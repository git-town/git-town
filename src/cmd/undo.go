package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	runState, err := runstate.Load(&run.Backend)
	if err != nil {
		return fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || runState.IsUnfinished() {
		return fmt.Errorf("nothing to undo")
	}
	undoRunState := runState.CreateUndoRunState()
	return runstate.Execute(&undoRunState, &run, nil)
}
