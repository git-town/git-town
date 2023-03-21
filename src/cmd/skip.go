package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

const skipDesc = "Restarts the last run git-town command by skipping the current branch"

func skipCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "skip",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured),
		Short:   skipDesc,
		Long:    long(skipDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return skip(repo)
		},
	}
}

func skip(repo *git.ProdRepo) error {
	runState, err := runstate.Load(repo)
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
	return runstate.Execute(&skipRunState, repo, nil)
}
