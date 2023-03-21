package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const syncStrategyDesc = "Displays or sets your sync strategy"

const syncStrategyHelp = `
The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`

func syncStrategyCommand(repo *git.ProdRepo) *cobra.Command {
	var globalFlag bool
	syncStrategyCmd := cobra.Command{
		Use:     "sync-strategy [(merge | rebase)]",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: ensure(repo, isRepository),
		Short:   syncStrategyDesc,
		Long:    long(syncStrategyDesc, syncStrategyHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return configSyncStrategy(args, globalFlag, repo)
		},
	}
	syncStrategyCmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global sync strategy")
	return &syncStrategyCmd
}

func configSyncStrategy(args []string, globalFlag bool, repo *git.ProdRepo) error {
	if len(args) > 0 {
		return setSyncStrategy(globalFlag, repo, args[0])
	}
	return printSyncStrategy(globalFlag, repo)
}

func printSyncStrategy(globalFlag bool, repo *git.ProdRepo) error {
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

func setSyncStrategy(globalFlag bool, repo *git.ProdRepo, value string) error {
	syncStrategy, err := config.ToSyncStrategy(value)
	if err != nil {
		return err
	}
	if globalFlag {
		return repo.Config.SetSyncStrategyGlobal(syncStrategy)
	}
	return repo.Config.SetSyncStrategy(syncStrategy)
}
