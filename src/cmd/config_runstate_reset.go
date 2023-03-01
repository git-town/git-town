package cmd

import (
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
		Short: "Resets the persisted runstate",
		Run: func(cmd *cobra.Command, args []string) {
			filepath, err := runstate.PersistenceFilename(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot determine the runstate file: %w", err))
			}
			err = os.Remove(filepath)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
	}
}
