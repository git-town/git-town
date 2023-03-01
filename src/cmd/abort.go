package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func abortCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "abort",
		Short: "Aborts the last run git-town command",
		Run: func(cmd *cobra.Command, args []string) {
			runState, err := runstate.Load(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot load previous run state: %w", err))
			}
			if runState == nil || !runState.IsUnfinished() {
				cli.Exit(fmt.Errorf("nothing to abort"))
			}
			abortRunState := runState.CreateAbortRunState()
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintConnectorAction)
			if err != nil {
				cli.Exit(err)
			}
			err = runstate.Execute(&abortRunState, repo, connector)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
		GroupID: "errors",
	}
}
