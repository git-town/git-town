package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

func resetConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Resets your Git Town configuration",
		Run: func(cmd *cobra.Command, args []string) {
			err := prodRepo.Config.RemoveLocalGitConfiguration()
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(prodRepo)
		},
	}
}
