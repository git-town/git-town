package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		command.PrintLabelAndValue("Main branch", git.GetPrintableMainBranch())
		command.PrintLabelAndValue("Perennial branches", git.GetPrintablePerennialBranchTrees())
		mainBranch := git.Config().GetMainBranch()
		if mainBranch != "" {
			command.PrintLabelAndValue("Branch Ancestry", git.GetPrintableBranchTree(mainBranch))
		}
		command.PrintLabelAndValue("Pull branch strategy", git.Config().GetPullBranchStrategy())
		command.PrintLabelAndValue("New Branch Push Flag", git.GetPrintableNewBranchPushFlag())
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
		prompt.ConfigurePerennialBranches()
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
