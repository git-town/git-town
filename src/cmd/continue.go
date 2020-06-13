package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Restarts the last run git-town command after having resolved conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState(prodRepo)
		if err != nil {
			fmt.Printf("cannot load previous run state: %v\n", err)
			os.Exit(1)
		}
		if runState == nil || !runState.IsUnfinished() {
			cli.Exit("Nothing to continue")
		}
		hasConflicts, err := prodRepo.Silent.HasConflicts()
		if err != nil {
			cli.Exit(err)
		}
		if hasConflicts {
			fmt.Println("Error: you must resolve the conflicts before continuing")
			os.Exit(1)
		}
		err = steps.Run(runState, prodRepo, drivers.Load(prodRepo.Configuration, &prodRepo.Silent))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func init() {
	RootCmd.AddCommand(continueCmd)
}
