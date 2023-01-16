package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func addAbortCmd(rootCmd *cobra.Command) {
	abortCmd := &cobra.Command{
		Use:   "abort",
		Short: "Aborts the last run git-town command",
		Run: func(cmd *cobra.Command, args []string) {
			runState, err := runstate.Load(prodRepo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot load previous run state: %w", err))
			}
			if runState == nil || !runState.IsUnfinished() {
				cli.Exit(fmt.Errorf("nothing to abort"))
			}
			abortRunState := runState.CreateAbortRunState()
			err = runstate.Execute(&abortRunState, prodRepo, hosting.NewDriver(&prodRepo.Config, &prodRepo.Silent, cli.PrintDriverAction))
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
	rootCmd.AddCommand(abortCmd)
}
