package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const configSyncStrategySummary = "Displays or sets your sync strategy"

const configSyncStrategyDesc = `
The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`

func syncStrategyCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "sync-strategy [(merge | rebase)]",
		Args:  cobra.MaximumNArgs(1),
		Short: configSyncStrategySummary,
		Long:  long(configSyncStrategySummary, configSyncStrategyDesc),
		RunE:  runConfigSyncStrategy,
	}
	addDebugFlag(&cmd)
	cmd.Flags().Bool(globalFlagName, false, "When set, displays or sets the sync strategy for all repos on this machine")
	return &cmd
}

func runConfigSyncStrategy(cmd *cobra.Command, args []string) error {
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
	globalFlag := readBoolFlag(cmd, globalFlagName)
	if len(args) > 0 {
		return setSyncStrategy(globalFlag, &repo, args[0])
	}
	return printSyncStrategy(globalFlag, &repo)
}

func printSyncStrategy(globalFlag bool, repo *git.PublicRepo) error {
	var strategy config.SyncStrategy
	var err error
	if globalFlag {
		strategy, err = repo.Config.SyncStrategyGlobal()
	} else {
		strategy, err = repo.Config.SyncStrategy()
	}
	if err != nil {
		return err
	}
	cli.Println(strategy)
	return nil
}

func setSyncStrategy(globalFlag bool, repo *git.PublicRepo, value string) error {
	syncStrategy, err := config.ToSyncStrategy(value)
	if err != nil {
		return err
	}
	if globalFlag {
		return repo.Config.SetSyncStrategyGlobal(syncStrategy)
	}
	return repo.Config.SetSyncStrategy(syncStrategy)
}
