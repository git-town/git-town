package cmd

import (
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/spf13/cobra"
)

const resetConfigDesc = "Resets your Git Town configuration"

func resetConfigCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: resetConfigDesc,
		Long:  long(resetConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetStatus(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func resetStatus(debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: true,
	})
	if err != nil {
		return err
	}
	_, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	return run.Config.RemoveLocalGitConfiguration()
}
