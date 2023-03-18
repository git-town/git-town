package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func undoCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   "Undoes the last run git-town command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUndo(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runUndo(debug bool) error {
	repo, err := LoadPublicRepo(RepoArgs{
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
	if runState == nil || runState.IsUnfinished() {
		return fmt.Errorf("nothing to undo")
	}
	undoRunState := runState.CreateUndoRunState()
	return runstate.Execute(&undoRunState, &repo, nil)
}
