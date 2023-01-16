package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func continueCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "continue",
		Short: "Restarts the last run git-town command after having resolved conflicts",
		Run: func(cmd *cobra.Command, args []string) {
			runState, err := runstate.Load(prodRepo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot load previous run state: %w", err))
			}
			if runState == nil || !runState.IsUnfinished() {
				cli.Exit(fmt.Errorf("nothing to continue"))
			}
			hasConflicts, err := prodRepo.Silent.HasConflicts()
			if err != nil {
				cli.Exit(err)
			}
			if hasConflicts {
				cli.Exit(fmt.Errorf("you must resolve the conflicts before continuing"))
			}
			err = runstate.Execute(runState, prodRepo, hosting.NewDriver(&prodRepo.Config, &prodRepo.Silent, cli.PrintDriverAction))
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
}
