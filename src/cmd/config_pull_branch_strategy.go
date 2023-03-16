package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pullBranchStrategyCommand(repo *git.PublicRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "pull-branch-strategy [(rebase | merge)]",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: ensure(repo, isRepository),
		Short:   "Displays or sets your pull branch strategy",
		Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return setPullBranchStrategy(args[0], repo)
			}
			return displayPullBranchStrategy(repo)
		},
	}
}

func displayPullBranchStrategy(repo *git.PublicRepo) error {
	pullBranchStrategy, err := repo.Config.PullBranchStrategy()
	if err != nil {
		return err
	}
	cli.Println(pullBranchStrategy)
	return nil
}

func setPullBranchStrategy(value string, repo *git.PublicRepo) error {
	pullBranchStrategy, err := config.NewPullBranchStrategy(value)
	if err != nil {
		return err
	}
	return repo.Config.SetPullBranchStrategy(pullBranchStrategy)
}
