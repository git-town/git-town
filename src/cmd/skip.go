package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var skipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Restarts the last run git-town command by skipping the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		runState := steps.LoadPreviousRunState()
		if runState == nil || !runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to skip")
		}
		if !runState.UnfinishedDetails.CanSkip {
			util.ExitWithErrorMessage("Cannot skip branch that resulted in conflicts")
		}
		skipRunState := runState.CreateSkipRunState()
		steps.Run(&skipRunState)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(skipCmd)
}
