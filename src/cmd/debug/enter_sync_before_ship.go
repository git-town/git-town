package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/spf13/cobra"
)

func enterSyncBeforeShip() *cobra.Command {
	return &cobra.Command{
		Use: "sync-before-ship",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := enter.SyncBeforeShip(false, dialogTestInputs.Next())
			return err
		},
	}
}
