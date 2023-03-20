package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func skipCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   "Restarts the last run git-town command by skipping the current branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkip(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runSkip(debug bool) error {
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
		return fmt.Errorf("nothing to skip")
	}
	if !runState.UnfinishedDetails.CanSkip {
		return fmt.Errorf("cannot skip branch that resulted in conflicts")
	}
	skipRunState := runState.CreateSkipRunState()
	return runstate.Execute(&skipRunState, &repo, nil)
}
