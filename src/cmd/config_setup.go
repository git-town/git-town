package cmd

import (
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func setupConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Prompts to setup your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := dialog.ConfigureMainBranch(repo)
			if err != nil {
				return err
			}
			return dialog.ConfigurePerennialBranches(repo)
		},
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
	}
}
