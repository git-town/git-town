package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
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
			err := setMainBranch(args[0], prodRepo)
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
	cli.Println(git.GetPrintableMainBranch())
}

func setMainBranch(branchName string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Silent.HasLocalBranch(branchName)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branchName)
	}
	git.Config().SetMainBranch(branchName)
	return nil
}

func init() {
	RootCmd.AddCommand(mainBranchCommand)
}
