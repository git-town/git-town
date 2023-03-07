package cmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return displayPullBranchStrategy(repo)
			}
			return setPullBranchStrategy(args[0], repo)
		},
		ValidArgs: []string{"merge", "rebase"},
		Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}

func displayPullBranchStrategy(repo *git.ProdRepo) error {
	pullBranchStrategy, err := repo.Config.PullBranchStrategy()
	if err != nil {
		return err
	}
	cli.Println(pullBranchStrategy)
	return nil
}

func setPullBranchStrategy(value string, repo *git.ProdRepo) error {
	pullBranchStrategy, err := config.ToPullBranchStrategy(value)
	if err != nil {
		return err
	}
	return repo.Config.SetPullBranchStrategy(pullBranchStrategy)
}
