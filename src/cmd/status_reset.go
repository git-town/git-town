package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository),
		Short:   statusResetDesc,
		Long:    long(statusResetDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusReset(repo)
		},
	}
}

func statusReset(repo *git.ProdRepo) error {
	err := runstate.Delete(repo)
	if err != nil {
		return err
	}
	fmt.Println("Runstate file deleted.")
	return nil
}
