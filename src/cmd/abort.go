package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/steps"

	"github.com/spf13/cobra"
)

var abortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Aborts the last run git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState(prodRepo)
		if err != nil {
			cli.Exit(fmt.Errorf("cannot load previous run state: %w", err))
		}
		if runState == nil || !runState.IsUnfinished() {
			cli.Exit(fmt.Errorf("nothing to abort"))
		}
		abortRunState := runState.CreateAbortRunState()
		err = steps.Run(&abortRunState, prodRepo, drivers.Load(prodRepo.Config, &prodRepo.Silent, cli.PrintDriverAction))
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
	RootCmd.AddCommand(abortCmd)
}
