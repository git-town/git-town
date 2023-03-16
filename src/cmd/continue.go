package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func continueCmd(repo *git.PublicRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "continue",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured),
		Short:   "Restarts the last run git-town command after having resolved conflicts",
		RunE: func(cmd *cobra.Command, args []string) error {
			runState, err := runstate.Load(repo)
			if err != nil {
				return fmt.Errorf("cannot load previous run state: %w", err)
			}
			if runState == nil || !runState.IsUnfinished() {
				return fmt.Errorf("nothing to continue")
			}
			hasConflicts, err := repo.HasConflicts()
			if err != nil {
				return err
			}
			if hasConflicts {
				return fmt.Errorf("you must resolve the conflicts before continuing")
			}
			connector, err := hosting.NewConnector(repo.Config, &repo.InternalRepo, cli.PrintConnectorAction)
			if err != nil {
				return err
			}
			return runstate.Execute(runState, repo, connector)
		},
	}
}
