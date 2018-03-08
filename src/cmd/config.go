package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		printConfig()
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
		resetConfig()
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
		setupConfig()
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

func printConfig() {
	fmt.Println()
	util.PrintLabelAndValue("Main branch", git.GetPrintableMainBranch())
	util.PrintLabelAndValue("Perennial branches", git.GetPrintablePerennialBranches())

	mainBranch := git.GetMainBranch()
	if mainBranch != "" {
		util.PrintLabelAndValue("Branch Ancestry", git.GetPrintableBranchTree(mainBranch))
	}

	util.PrintLabelAndValue("Pull branch strategy", git.GetPullBranchStrategy())
	util.PrintLabelAndValue("New Branch Push Flag", git.GetPrintableNewBranchPushFlag())
}

func resetConfig() {
	git.RemoveAllConfiguration()
}

func setupConfig() {
	prompt.ConfigureMainBranch()
	prompt.ConfigurePerennialBranches()
}

func init() {
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	RootCmd.AddCommand(configCommand)
}
