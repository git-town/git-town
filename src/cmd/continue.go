package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Restarts the last run git-town command after having resolved conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		runState := steps.LoadPreviousRunState()
		if runState == nil || !runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to continue")
		}
		git.EnsureDoesNotHaveConflicts()
		steps.Run(runState)
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
	RootCmd.AddCommand(continueCmd)
}
