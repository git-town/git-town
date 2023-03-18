package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func resetRunstateCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: "Resets the current suspended Git Town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusReset(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runStatusReset(debug bool) error {
	repo, err := LoadRepo(RepoArgs{
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	err = runstate.Delete(&repo)
	if err != nil {
		return err
	}
	fmt.Println("Runstate file deleted.")
	return nil
}
