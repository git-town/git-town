package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cfmt"
	"github.com/git-town/git-town/src/git"
	"github.com/spf13/cobra"
)

var mainBranchCommand = &cobra.Command{
	Use:   "main-branch [<branch>]",
	Short: "Displays or sets your main development branch",
	Long: `Displays or sets your main development branch

The main branch is the Git branch from which new feature branches are cut.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printMainBranch()
		} else {
			err := setMainBranch(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

func printMainBranch() {
	cfmt.Println(git.GetPrintableMainBranch())
}

func setMainBranch(branchName string) error {
	if !git.HasBranch(branchName) {
		return fmt.Errorf("there is no branch named %q", branchName)
	}
	git.Config().SetMainBranch(branchName)
	return nil
}

func init() {
	RootCmd.AddCommand(mainBranchCommand)
}
