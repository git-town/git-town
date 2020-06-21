package cmd

import (
	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var abortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Aborts the last run git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState(prodRepo)
		if err != nil {
			cli.Exit("cannot load previous run state: %v\n", err)
		}
		if runState == nil || !runState.IsUnfinished() {
			cli.Exit("Nothing to abort")
		}
		abortRunState := runState.CreateAbortRunState()
		err = steps.Run(&abortRunState, prodRepo, drivers.Load(prodRepo.Config, &prodRepo.Silent, cli.PrintLog))
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
