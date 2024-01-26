package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncFeatureStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-feature-strategy",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialog.LoadTestInputs(os.Environ())
			_, _, err := enter.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, dialogTestInputs.Next())
			return err
		},
	}
}
