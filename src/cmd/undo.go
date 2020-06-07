package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"

	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the last run git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState()
		if err != nil {
			fmt.Printf("cannot load previous run state: %v\n", err)
			os.Exit(1)
		}
		if runState == nil || runState.IsUnfinished() {
			cli.Exit("Nothing to undo")
		}
		undoRunState := runState.CreateUndoRunState()
		err = steps.Run(&undoRunState, repo(), nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := git.ValidateIsRepository(); err != nil {
			return err
		}
		return validateIsConfigured(repo())
	},
}

func init() {
	RootCmd.AddCommand(undoCmd)
}
