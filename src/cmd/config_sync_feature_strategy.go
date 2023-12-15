package cmd

import (
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/spf13/cobra"
)

const syncFeatureStrategyDesc = "Displays or sets your sync-feature strategy"

const syncFeatureStrategyHelp = `
The sync-feature strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`

func syncFeatureStrategyCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "When set, displays or sets the sync-feature strategy for all repos on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "sync-feature-strategy [(merge | rebase)]",
		Args:  cobra.MaximumNArgs(1),
		Short: syncFeatureStrategyDesc,
		Long:  long(syncFeatureStrategyDesc, syncFeatureStrategyHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigSyncFeatureStrategy(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeConfigSyncFeatureStrategy(args []string, global, verbose bool) error {
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
		return setSyncFeatureStrategy(global, repo.Runner, args[0])
	}
	return printSyncFeatureStrategy(global, repo.Runner)
}

func printSyncFeatureStrategy(globalFlag bool, run *git.ProdRunner) error {
	var strategy configdomain.SyncFeatureStrategy
	var err error
	if globalFlag {
		strategy, err = run.GitTown.SyncFeatureStrategyGlobal()
	} else {
		strategy, err = run.GitTown.SyncFeatureStrategy()
	}
	if err != nil {
		return err
	}
	io.Println(strategy)
	return nil
}

func setSyncFeatureStrategy(globalFlag bool, run *git.ProdRunner, value string) error {
	syncFeatureStrategy, err := configdomain.NewSyncFeatureStrategy(value)
	if err != nil {
		return err
	}
	if globalFlag {
		return run.GitTown.SetSyncFeatureStrategyGlobal(syncFeatureStrategy)
	}
	return run.GitTown.SetSyncFeatureStrategy(syncFeatureStrategy)
}
