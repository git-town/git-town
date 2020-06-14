package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var skipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Restarts the last run git-town command by skipping the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState(prodRepo)
		if err != nil {
			cli.Exit(fmt.Errorf("cannot load previous run state: %v", err))
		}
		if runState == nil || !runState.IsUnfinished() {
			cli.Exit("Nothing to skip")
		}
		if !runState.UnfinishedDetails.CanSkip {
			cli.Exit("Cannot skip branch that resulted in conflicts")
		}
		skipRunState := runState.CreateSkipRunState()
		err = steps.Run(&skipRunState, prodRepo, nil)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func init() {
	RootCmd.AddCommand(skipCmd)
}
