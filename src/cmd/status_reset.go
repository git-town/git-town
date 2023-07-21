package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  long(statusResetDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusReset(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func statusReset(debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsConfigured:  false,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	err = runstate.Delete(&run.Backend)
	if err != nil {
		return err
	}
	fmt.Println("Runstate file deleted.")
	return nil
}
