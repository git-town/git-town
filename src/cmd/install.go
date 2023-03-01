package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func installCommand(repo *git.ProdRepo, rootCmd *cobra.Command) *cobra.Command {
	installCmd := cobra.Command{
		Use:     "install",
		Short:   "Commands to set up Git Town on your computer",
		Args:    cobra.NoArgs,
		GroupID: "setup",
	}
	installCmd.AddCommand(aliasCommand(repo))
	installCmd.AddCommand(completionsCmd(rootCmd))
	return &installCmd
}
