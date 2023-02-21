package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func proposalsCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "proposals",
		Short: "Displays currently active proposals and their status",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Println("hello")
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			return repo.Config.ValidateIsOnline()
		},
	}
}
