package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func continueCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "continue",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   "Restarts the last run git-town command after having resolved conflicts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runContinue(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runContinue(debug bool) error {
	repo, err := LoadRepo(RepoArgs{
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	runState, err := runstate.Load(&repo)
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
	return runstate.Execute(runState, &repo, connector)
}
