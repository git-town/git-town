package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func setupConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "setup",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
		Short:   "Prompts to setup your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := validate.ConfigureMainBranch(repo)
			if err != nil {
				return err
			}
			return validate.ConfigurePerennialBranches(repo)
		},
	}
}
