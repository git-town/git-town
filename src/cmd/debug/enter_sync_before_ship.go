package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/spf13/cobra"
)

func enterSyncBeforeShip() *cobra.Command {
	return &cobra.Command{
		Use: "sync-before-ship",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterSyncBeforeShip(false, dialogTestInputs.Next())
			return err
		},
	}
}
