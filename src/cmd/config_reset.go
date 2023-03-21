package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const configResetDesc = "Resets your Git Town configuration"

func resetConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
		Short:   configResetDesc,
		Long:    long(completionsDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo.Config.RemoveLocalGitConfiguration()
		},
	}
}
