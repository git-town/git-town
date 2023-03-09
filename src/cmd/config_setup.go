package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func setupConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Prompts to setup your Git Town configuration",
		Run: func(cmd *cobra.Command, args []string) {
			err := dialog.ConfigureMainBranch(repo)
			if err != nil {
				cli.Exit(err)
			}
			err = dialog.ConfigurePerennialBranches(repo)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
	}
}
