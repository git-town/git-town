package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/spf13/cobra"
)

func enterSyncUpstream() *cobra.Command {
	return &cobra.Command{
		Use: "sync-upstream",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := dialog.EnterSyncUpstream(true, dialogTestInputs.Next())
			return err
		},
	}
}
