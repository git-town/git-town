package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Continues the previous git-town command that encountered conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		runState := steps.LoadPreviousRunState()
		if runState == nil || !runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to continue")
		}
		if skipFlag {
			if runState.UnfinishedDetails.CanSkip {
				skipRunState := runState.CreateSkipRunState()
				steps.Run(&skipRunState)
			} else {
				util.ExitWithErrorMessage("Cannot skip branch that resulted in conflicts")
			}
		} else {
			git.EnsureDoesNotHaveConflicts()
			steps.Run(runState)
		}
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
	continueCmd.Flags().BoolVar(&skipFlag, "skip", false, "Skip the branch that resulted in conflicts")
	RootCmd.AddCommand(continueCmd)
}
