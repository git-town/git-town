package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/spf13/cobra"
)

var mainBranchCommand = &cobra.Command{
	Use:   "main-branch [<branch>]",
	Short: "Displays or sets your main branch",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		if len(args) == 0 {
			printMainBranch()
		} else {
			setMainBranch(args[0])
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func printMainBranch() {
	fmt.Println(git.GetPrintableMainBranch())
}

func setMainBranch(branchName string) {
	git.EnsureHasBranch(branchName)
	git.SetMainBranch(branchName)
}

func init() {
	RootCmd.AddCommand(mainBranchCommand)
}
