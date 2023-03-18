package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pullBranchStrategyCommand() *cobra.Command {
	debug := false
	cmd := &cobra.Command{
		Use:   "pull-branch-strategy [(rebase | merge)]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Displays or sets your pull branch strategy",
		Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigurePullBranchStrategy(debug, args)
		},
	}
	debugFlag(cmd, &debug)
	return cmd
}

func runConfigurePullBranchStrategy(debug bool, args []string) error {
	repo, err := LoadRepo(RepoArgs{
		omitBranchNames:      true,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return setPullBranchStrategy(args[0], &repo)
	}
	return displayPullBranchStrategy(&repo)
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
