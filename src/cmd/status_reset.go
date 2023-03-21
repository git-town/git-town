package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

const statusResetSummary = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetSummary,
		Long:  long(statusResetSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusReset(debug)
		},
	}
	debugFlagOld(&cmd, &debug)
	return &cmd
}

func runStatusReset(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	err = runstate.Delete(&repo)
	if err != nil {
		return err
	}
	fmt.Println("Runstate file deleted.")
	return nil
}
