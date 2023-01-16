package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func mainbranchConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "main-branch [<branch>]",
		Short: "Displays or sets your main development branch",
		Long: `Displays or sets your main development branch

The main branch is the Git branch from which new feature branches are cut.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				printMainBranch(prodRepo)
			} else {
				err := setMainBranch(args[0], prodRepo)
				if err != nil {
					cli.Exit(err)
				}
			}
		},
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(prodRepo)
		},
	}
}

func printMainBranch(repo *git.ProdRepo) {
	cli.Println(cli.StringSetting(repo.Config.MainBranch()))
}

func setMainBranch(branchName string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Silent.HasLocalBranch(branchName)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branchName)
	}
	return repo.Config.SetMainBranch(branchName)
}
