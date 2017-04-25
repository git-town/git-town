package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/util"
	"github.com/spf13/cobra"
)

var resetFlag bool
var setupFlag bool

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays or updates your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if resetFlag {
			resetConfig()
		} else if setupFlag {
			setupConfig()
		} else {
			printConfig()
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
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
	util.PrintLabelAndValue("git-hack push flag", git.GetPrintableHackPushFlag())
}

func resetConfig() {
	git.RemoveAllConfiguration()
}

func setupConfig() {
	prompt.ConfigureMainBranch()
	prompt.ConfigurePerennialBranches()
}

func init() {
	configCommand.Flags().BoolVar(&resetFlag, "reset", false, "Remove all Git Town configuration from the current repository")
	configCommand.Flags().BoolVar(&setupFlag, "setup", false, "Run the Git Town configuration wizard")
	RootCmd.AddCommand(configCommand)
}
