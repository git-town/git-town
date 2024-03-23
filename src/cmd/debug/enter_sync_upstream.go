package debug

import (
	"os"

	"github.com/git-town/git-town/v13/src/cli/dialog"
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterSyncUpstream() *cobra.Command {
	return &cobra.Command{
		Use: "sync-upstream",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncUpstream(true, dialogTestInputs.Next())
			return err
		},
	}
}
