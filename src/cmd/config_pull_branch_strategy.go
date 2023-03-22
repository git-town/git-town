package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const pullBranchStrategyDesc = "Displays or sets your pull branch strategy"

const pullBranchStrategyHelp = `
The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`

func pullBranchStrategyCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:   "pull-branch-strategy [(rebase | merge)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pullBranchStrategyDesc,
		Long:  long(pullBranchStrategyDesc, pullBranchStrategyHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigurePullBranchStrategy(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runConfigurePullBranchStrategy(args []string, debug bool) error {
	repo, exit, err := LoadProdRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 debug,
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

func displayPullBranchStrategy(repo *git.ProdRepo) error {
	pullBranchStrategy, err := repo.Config.PullBranchStrategy()
	if err != nil {
		return err
	}
	cli.Println(pullBranchStrategy)
	return nil
}

func setPullBranchStrategy(value string, repo *git.ProdRepo) error {
	pullBranchStrategy, err := config.NewPullBranchStrategy(value)
	if err != nil {
		return err
	}
	return repo.Config.SetPullBranchStrategy(pullBranchStrategy)
}
