package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pullBranchStrategyCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "pull-branch-strategy [(rebase | merge)]",
		Short: "Displays or sets your pull branch strategy",
		Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				displayPullBranchStrategy(repo)
			} else {
				err := updatePullBranchStrategy(args[0], repo)
				if err != nil {
					cli.Exit(err)
				}
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 && args[0] != "rebase" && args[0] != "merge" {
				return fmt.Errorf("invalid value: %q", args[0])
			}
			return cobra.MaximumNArgs(1)(cmd, args)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}

func displayPullBranchStrategy(repo *git.ProdRepo) {
	cli.Println(repo.Config.PullBranchStrategy())
}

func updatePullBranchStrategy(value string, repo *git.ProdRepo) error {
	pullBranchStrategy, err := config.NewPullBranchStrategy(value)
	if err != nil {
		return err
	}
	return repo.Config.SetPullBranchStrategy(pullBranchStrategy)
}
