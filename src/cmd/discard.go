package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/steps"
	"github.com/git-town/git-town/src/util"

	"github.com/spf13/cobra"
)

var discardCmd = &cobra.Command{
	Use:   "discard",
	Short: "Discards the saved state of the previous git-town command",
	Run: func(cmd *cobra.Command, args []string) {
		err := steps.DeletePreviousRunState()
		if err != nil {
			fmt.Printf("cannot delete previous run state: %v", err)
			os.Exit(1)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(discardCmd)
}
