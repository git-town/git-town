package cmd

import (
	"github.com/git-town/git-town/v8/src/cli"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/flags"
	"github.com/git-town/git-town/v8/src/git"
	"github.com/spf13/cobra"
)

const pullBranchDesc = "Displays or sets your pull branch strategy"

const pullBranchHelp = `
The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`

func pullBranchStrategyCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "pull-branch-strategy [(rebase | merge)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pullBranchDesc,
		Long:  long(pullBranchDesc, pullBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pullBranchStrategy(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func pullBranchStrategy(args []string, debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames:       true,
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setPullBranchStrategy(args[0], &run)
	}
	return displayPullBranchStrategy(&run)
}

func displayPullBranchStrategy(run *git.ProdRunner) error {
	pullBranchStrategy, err := run.Config.PullBranchStrategy()
	if err != nil {
		return err
	}
	cli.Println(pullBranchStrategy)
	return nil
}

func setPullBranchStrategy(value string, run *git.ProdRunner) error {
	pullBranchStrategy, err := config.NewPullBranchStrategy(value)
	if err != nil {
		return err
	}
	return run.Config.SetPullBranchStrategy(pullBranchStrategy)
}
