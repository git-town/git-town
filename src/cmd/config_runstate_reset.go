package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func resetRunstateCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets the current interrupted Git Town command",
		Run: func(cmd *cobra.Command, args []string) {
			filepath, err := runstate.PersistenceFilename(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot determine the runstate file: %w", err))
			}
			err = os.Remove(filepath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					fmt.Println("Runstate doesn't exist.")
					return
				}
				cli.Exit(fmt.Errorf("cannot delete runstate file: %w", err))
			}
			fmt.Println("Runstate file deleted.")
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}
