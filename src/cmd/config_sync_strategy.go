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

func syncStrategyCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	addGlobalFlag, readGlobalFlag := boolFlag("global", "When set, displays or sets the sync strategy for all repos on this machine")
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
	repo, exit, err := LoadPublicRepo(RepoArgs{
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
		return setSyncStrategy(global, &repo, args[0])
	}
	return printSyncStrategy(global, &repo)
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
