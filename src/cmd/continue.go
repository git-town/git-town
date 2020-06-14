package cmd

import (
	"fmt"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Restarts the last run git-town command after having resolved conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState(prodRepo)
		if err != nil {
			cli.Exit(fmt.Errorf("cannot load previous run state: %v", err))
		}
		if runState == nil || !runState.IsUnfinished() {
			cli.Exit("Nothing to continue")
		}
		hasConflicts, err := prodRepo.Silent.HasConflicts()
		if err != nil {
			cli.Exit(err)
		}
		if hasConflicts {
			cli.Exit(fmt.Errorf("you must resolve the conflicts before continuing"))
		}
		err = steps.Run(runState, prodRepo, drivers.Load(prodRepo.Configuration, &prodRepo.Silent))
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
	RootCmd.AddCommand(continueCmd)
}
