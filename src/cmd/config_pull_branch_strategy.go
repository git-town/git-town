package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
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
				displayPullBranchStrategy(repo)
			} else {
				setPullBranchStrategy(args[0], repo)
			}
			return nil
		},
		ValidArgs: []string{"merge", "rebase"},
		Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
}

func displayPullBranchStrategy(repo *git.ProdRepo) {
	cli.Println(repo.Config.PullBranchStrategy())
}

func setPullBranchStrategy(value string, repo *git.ProdRepo) error {
	return repo.Config.SetPullBranchStrategy(value)
}
