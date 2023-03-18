package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func resetRunstateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository),
		Short:   "Resets the current suspended Git Town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runstate.Delete(repo)
			if err != nil {
				return err
			}
			fmt.Println("Runstate file deleted.")
			return nil
		},
	}
}
