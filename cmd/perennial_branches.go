package cmd

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/spf13/cobra"
)

var branchToAdd string
var branchToRemove string

var perennialBranchesCommand = &cobra.Command{
	Use:   "perennial-branches",
	Short: "Displays or updates your perennial branches",
	Run: func(cmd *cobra.Command, args []string) {
		if branchToAdd != "" {
			git.EnsureHasBranch(branchToAdd)
			git.EnsureIsNotMainBranch(branchToAdd, fmt.Sprintf("'%s' is already set as the main branch", branchToAdd))
			git.EnsureIsNotPerennialBranch(branchToAdd, fmt.Sprintf("'%s' is already a perennial branch", branchToAdd))
			git.AddToPerennialBranches(branchToAdd)
			return
		}

		if branchToRemove != "" {
			git.EnsureIsPerennialBranch(branchToRemove, fmt.Sprintf("'%s' is not a perennial branch", branchToRemove))
			git.RemoveFromPerennialBranches(branchToRemove)
			return
		}

		output := strings.Join(git.GetPerennialBranches(), "\n")
		if output == "" {
			output = "[none]"
		}
		fmt.Println(output)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func init() {
	perennialBranchesCommand.Flags().StringVar(&branchToAdd, "add", "", "Specify a branch to add")
	perennialBranchesCommand.Flags().StringVar(&branchToRemove, "remove", "", "Specify a branch to remove")
	RootCmd.AddCommand(perennialBranchesCommand)
}
