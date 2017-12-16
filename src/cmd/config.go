package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

var resetFlag bool
var setupFlag bool

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays or resets your Git Town configuration",
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
		return util.FirstError(
			validateMaxArgsFunc(args, 0),
			git.ValidateIsRepository,
		)
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
	util.PrintLabelAndValue("git-hack push flag", git.GetPrintableNewBranchPushFlag())
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
