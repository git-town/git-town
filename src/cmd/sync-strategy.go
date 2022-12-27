package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

var syncStrategyCommand = &cobra.Command{
	Use:   "sync-strategy [(merge | rebase)]",
	Short: "Displays or sets your sync strategy",
	Long: `Displays or sets your sync strategy

The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(prodRepo.Config.SyncStrategy())
		} else {
			err := prodRepo.Config.SetSyncStrategy(args[0])
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
		return ValidateIsRepository(prodRepo)
	},
}

func init() {
	RootCmd.AddCommand(syncStrategyCommand)
}
