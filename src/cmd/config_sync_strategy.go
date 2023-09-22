package cmd

import (
	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/spf13/cobra"
)

const syncStrategyDesc = "Displays or sets your sync strategy"

const syncStrategyHelp = `
The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`

func syncStrategyCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "When set, displays or sets the sync strategy for all repos on this machine")
	cmd := cobra.Command{
		Use:   "sync-strategy [(merge | rebase)]",
		Args:  cobra.MaximumNArgs(1),
		Short: syncStrategyDesc,
		Long:  long(syncStrategyDesc, syncStrategyHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSyncStrategy(args, readGlobalFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func runConfigSyncStrategy(args []string, global, debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
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
	cli.Println(strategy)
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
