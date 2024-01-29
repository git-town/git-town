package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs"
	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/spf13/cobra"
)

func enterSyncUpstream() *cobra.Command {
	return &cobra.Command{
		Use: "sync-upstream",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.SyncUpstream(true, dialogTestInputs.Next())
			return err
		},
	}
}
