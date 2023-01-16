package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func discardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "discard",
		Short: "Discards the saved state of the previous git-town command",
		Run: func(cmd *cobra.Command, args []string) {
			err := runstate.Delete(prodRepo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot delete previous run state: %w", err))
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(prodRepo); err != nil {
				return err
			}
			return validateIsConfigured(prodRepo)
		},
	}
}
