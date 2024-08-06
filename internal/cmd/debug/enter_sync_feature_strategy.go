package debug

import (
	"os"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncFeatureStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-feature-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, dialogTestInputs.Next())
			return err
		},
	}
}
