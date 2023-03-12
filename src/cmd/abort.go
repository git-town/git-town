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
		Use:     "abort",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured),
		Short:   "Aborts the last run git-town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			runState, err := runstate.Load(repo)
			if err != nil {
				return fmt.Errorf("cannot load previous run state: %w", err)
			}
			if runState == nil || !runState.IsUnfinished() {
				return fmt.Errorf("nothing to abort")
			}
			abortRunState := runState.CreateAbortRunState()
			connector, err := hosting.NewConnector(repo.Config, &repo.Silent, cli.PrintConnectorAction)
			if err != nil {
				return err
			}
			return runstate.Execute(&abortRunState, repo, connector)
		},
	}
}
