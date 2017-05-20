package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/spf13/cobra"
)

var pullBranchStrategyCommand = &cobra.Command{
	Use:   "pull-branch-strategy [(rebase | merge)]",
	Short: "Displays or sets your pull branch strategy",
	Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		if len(args) == 0 {
			printPullBranchStrategy()
		} else {
			setPullBranchStrategy(args[0])
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "rebase" && args[0] != "merge" {
			return fmt.Errorf("Invalid value: '%s'", args[0])
		}
		return validateMaxArgs(args, 1)
	},
}

func printPullBranchStrategy() {
	fmt.Println(git.GetPullBranchStrategy())
}

func setPullBranchStrategy(value string) {
	git.SetPullBranchStrategy(value)
}

func init() {
	RootCmd.AddCommand(pullBranchStrategyCommand)
}
