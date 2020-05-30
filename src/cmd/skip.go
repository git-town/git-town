package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"

	"github.com/spf13/cobra"
)

var skipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Restarts the last run git-town command by skipping the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		runState, err := steps.LoadPreviousRunState()
		if err != nil {
			fmt.Printf("cannot load previous run state: %v\n", err)
			os.Exit(1)
		}
		if runState == nil || !runState.IsUnfinished() {
			util.ExitWithErrorMessage("Nothing to skip")
		}
		if !runState.UnfinishedDetails.CanSkip {
			util.ExitWithErrorMessage("Cannot skip branch that resulted in conflicts")
		}
		skipRunState := runState.CreateSkipRunState()
		err = steps.Run(&skipRunState, git.NewProdRepo())
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
		return validateIsConfigured()
	},
}

func init() {
	RootCmd.AddCommand(skipCmd)
}
