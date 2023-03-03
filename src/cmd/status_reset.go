package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func resetRunstateCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets the current suspended Git Town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runstate.Delete(repo)
			if err != nil {
				return err
			}
			fmt.Println("Runstate file deleted.")
			return nil
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}
