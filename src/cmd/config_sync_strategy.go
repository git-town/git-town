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
			return syncStrategy(args, readGlobalFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func syncStrategy(args []string, global, debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: true,
	})
	if err != nil {
		return err
	}
	_, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setSyncStrategy(global, &run, args[0])
	}
	return printSyncStrategy(global, &run)
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
