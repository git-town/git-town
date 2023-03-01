package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func undoCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "undo",
		Short: "Undoes the last run git-town command",
		Run: func(cmd *cobra.Command, args []string) {
			runState, err := runstate.Load(repo)
			if err != nil {
				cli.Exit(fmt.Errorf("cannot load previous run state: %w", err))
			}
			if runState == nil || runState.IsUnfinished() {
				cli.Exit(fmt.Errorf("nothing to undo"))
			}
			undoRunState := runState.CreateUndoRunState()
			err = runstate.Execute(&undoRunState, repo, nil)
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
