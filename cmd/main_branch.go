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
		if len(args) == 0 {
			output := git.GetMainBranch()
			if output == "" {
				output = "[none]"
			}
			fmt.Println(output)
			return
		}

		newMainBranch := args[0]
		git.EnsureHasBranch(newMainBranch)
		git.SetMainBranch(newMainBranch)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 1)
	},
}

func init() {
	RootCmd.AddCommand(mainBranchCommand)
}
