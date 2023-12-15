package cmd

import (
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/spf13/cobra"
)

const pullBranchDesc = "Displays or sets your sync-perennial strategy"

const pullBranchHelp = `
The sync-perennial strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`

func syncPerennialStrategyCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "sync-perennial-strategy [(rebase | merge)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pullBranchDesc,
		Long:  long(pullBranchDesc, pullBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigPullBranch(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigPullBranch(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  false,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return setSyncPerennialStrategy(args[0], repo.Runner)
	}
	return displaySyncPerennialStrategy(repo.Runner)
}

func displaySyncPerennialStrategy(run *git.ProdRunner) error {
	syncPerennialStrategy, err := run.GitTown.SyncPerennialStrategy()
	if err != nil {
		return err
	}
	io.Println(syncPerennialStrategy)
	return nil
}

func setSyncPerennialStrategy(value string, run *git.ProdRunner) error {
	syncPerennialStrategy, err := configdomain.NewSyncPerennialStrategy(value)
	if err != nil {
		return err
	}
	return run.GitTown.SetSyncPerennialStrategy(syncPerennialStrategy)
}
