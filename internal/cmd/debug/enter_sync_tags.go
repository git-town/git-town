package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/spf13/cobra"
)

func enterSyncTags() *cobra.Command {
	return &cobra.Command{
		Use: "sync-tags",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncTags(true, dialogTestInputs.Next())
			return err
		},
	}
}
