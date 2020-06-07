package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cli.PrintLabelAndValue("Main branch", git.GetPrintableMainBranch())
		cli.PrintLabelAndValue("Perennial branches", git.GetPrintablePerennialBranchTrees())
		mainBranch := git.Config().GetMainBranch()
		if mainBranch != "" {
			cli.PrintLabelAndValue("Branch Ancestry", git.GetPrintableBranchTree(mainBranch))
		}
		cli.PrintLabelAndValue("Pull branch strategy", git.Config().GetPullBranchStrategy())
		cli.PrintLabelAndValue("New Branch Push Flag", git.GetPrintableNewBranchPushFlag())
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

var resetConfigCommand = &cobra.Command{
	Use:   "reset",
	Short: "Resets your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		git.Config().RemoveLocalGitConfiguration()
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

var setupConfigCommand = &cobra.Command{
	Use:   "setup",
	Short: "Prompts to setup your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		prompt.ConfigureMainBranch()
		prompt.ConfigurePerennialBranches(repo())
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

func init() {
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	RootCmd.AddCommand(configCommand)
}
