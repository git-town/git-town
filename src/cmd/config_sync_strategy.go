package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func syncStrategyCommand(repo *git.ProdRepo) *cobra.Command {
	var globalFlag bool
	syncStrategyCmd := cobra.Command{
		Use:   "sync-strategy [(merge | rebase)]",
		Short: "Displays or sets your sync strategy",
		Long: `Displays or sets your sync strategy

The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				printSyncStrategy(globalFlag, repo)
			} else {
				err := setSyncStrategy(globalFlag, repo, args[0])
				if err != nil {
					cli.Exit(err)
				}
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 && args[0] != "merge" && args[0] != "rebase" {
				return fmt.Errorf("invalid value: %q", args[0])
			}
			return cobra.MaximumNArgs(1)(cmd, args)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(repo)
		},
	}
	syncStrategyCmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets the global sync strategy")
	return &syncStrategyCmd
}

func printSyncStrategy(globalFlag bool, repo *git.ProdRepo) {
	var strategy string
	if globalFlag {
		strategy = repo.Config.SyncStrategyGlobal()
	} else {
		strategy = repo.Config.SyncStrategy()
	}
	cli.Println(strategy)
}

func setSyncStrategy(globalFlag bool, repo *git.ProdRepo, value string) error {
	if globalFlag {
		return repo.Config.SetSyncStrategyGlobal(value)
	}
	return repo.Config.SetSyncStrategy(value)
}
