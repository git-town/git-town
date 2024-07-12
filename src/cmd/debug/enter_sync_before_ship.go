package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterSyncBeforeShip() *cobra.Command {
	return &cobra.Command{
		Use: "sync-before-ship",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncBeforeShip(false, dialogTestInputs.Value.Next())
			return err
		},
	}
}
