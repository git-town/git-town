package cmd

import (
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/flags"
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
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames:       true,
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	return run.Config.RemoveLocalGitConfiguration()
}
