package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func resetConfigCommand(repo *git.PublicRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "reset",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
		Short:   "Resets your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo.Config.RemoveLocalGitConfiguration()
		},
	}
}
