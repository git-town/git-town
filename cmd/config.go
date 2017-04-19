package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
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
	util.PrintTitle("Main branch:")
	fmt.Println(util.Indent(git.GetPrintableMainBranch(), 1))
	fmt.Println()

	util.PrintTitle("Perennial branches:")
	fmt.Println(util.Indent(git.GetPrintablePerennialBranches(), 1))
	fmt.Println()

	mainBranch := git.GetMainBranch()
	if mainBranch != "" {
		util.PrintTitle("Branch Ancestry:")
		fmt.Println(util.Indent(git.GetPrintableBranchTree(mainBranch), 1))
		fmt.Println()
	}

	util.PrintTitle("Pull branch strategy:")
	fmt.Println(util.Indent(git.GetPullBranchStrategy(), 1))
	fmt.Println()

	util.PrintTitle("git-hack push flag:")
	fmt.Println(util.Indent(git.GetPrintableHackPushFlag(), 1))
	fmt.Println()
}

func resetConfig() {

}

func setupConfig() {

}

func init() {
	configCommand.Flags().BoolVar(&resetFlag, "reset", false, "Remove all Git Town configuration from the current repository")
	configCommand.Flags().BoolVar(&setupFlag, "setup", false, "Run the Git Town configuration wizard")
	RootCmd.AddCommand(configCommand)
}
