package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs"
	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/spf13/cobra"
)

func enterSyncBeforeShip() *cobra.Command {
	return &cobra.Command{
		Use: "sync-before-ship",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.SyncBeforeShip(false, dialogTestInputs.Next())
			return err
		},
	}
}
