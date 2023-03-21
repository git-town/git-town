package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  long(statusResetDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusReset(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
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
