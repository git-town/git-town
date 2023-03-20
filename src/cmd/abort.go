package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func abortCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "abort",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   "Aborts the last run git-town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAbort(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runAbort(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	runState, err := runstate.Load(&repo)
	if err != nil {
		return fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return fmt.Errorf("nothing to abort")
	}
	abortRunState := runState.CreateAbortRunState()
	connector, err := hosting.NewConnector(repo.Config, &repo.InternalRepo, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	return runstate.Execute(&abortRunState, &repo, connector)
}
