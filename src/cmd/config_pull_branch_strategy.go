package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const configPullBranchSummary = "Displays or sets your pull branch strategy"

const configPullBranchDesc = `
The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`

func pullBranchStrategyCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "pull-branch-strategy [(rebase | merge)]",
		Args:  cobra.MaximumNArgs(1),
		Short: configPullBranchSummary,
		Long:  long(configPullBranchSummary, configPullBranchDesc),
		RunE:  runConfigurePullBranchStrategy,
	}
	debugFlagOld(&cmd, &debug)
	return &cmd
}

func runConfigurePullBranchStrategy(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
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
