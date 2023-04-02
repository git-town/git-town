package cmd

import (
	"github.com/git-town/git-town/v7/src/flags"
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
	run, exit, err := LoadProdRunner(loadArgs{
		omitBranchNames:       true,
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	return run.Config.RemoveLocalGitConfiguration()
}
