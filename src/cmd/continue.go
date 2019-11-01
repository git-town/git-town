package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/steps"
	"github.com/Originate/git-town/src/util"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Restarts the last run git-town command after having resolved conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState()
		if err != nil {
			fmt.Printf("cannot load previous run state: %v\n", err)
			os.Exit(1)
		}
		if runState == nil || !runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to continue")
		}
		git.EnsureDoesNotHaveConflicts()
		err = steps.Run(runState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
	RootCmd.AddCommand(continueCmd)
}
