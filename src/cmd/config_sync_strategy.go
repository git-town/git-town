package cmd

import (
	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/io"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/spf13/cobra"
)

const syncStrategyDesc = "Displays or sets your sync strategy"

const syncStrategyHelp = `
The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`

func syncStrategyCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "When set, displays or sets the sync strategy for all repos on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "sync-strategy [(merge | rebase)]",
		Args:  cobra.MaximumNArgs(1),
		Short: syncStrategyDesc,
		Long:  long(syncStrategyDesc, syncStrategyHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigSyncStrategy(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeConfigSyncStrategy(args []string, global, verbose bool) error {
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
		return setSyncStrategy(global, &repo.Runner, args[0])
	}
	return printSyncStrategy(global, &repo.Runner)
}

func printSyncStrategy(globalFlag bool, run *git.ProdRunner) error {
	var strategy config.SyncStrategy
	var err error
	if globalFlag {
		strategy, err = run.Config.SyncStrategyGlobal()
	} else {
		strategy, err = run.Config.SyncStrategy()
	}
	if err != nil {
		return err
	}
	io.Println(strategy)
	return nil
}

func setSyncStrategy(globalFlag bool, run *git.ProdRunner, value string) error {
	syncStrategy, err := config.ToSyncStrategy(value)
	if err != nil {
		return err
	}
	if globalFlag {
		return run.Config.SetSyncStrategyGlobal(syncStrategy)
	}
	return run.Config.SetSyncStrategy(syncStrategy)
}
