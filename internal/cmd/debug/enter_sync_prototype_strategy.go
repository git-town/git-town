package debug

import (
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncPrototypeStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-prototype-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncPrototypeStrategy(configdomain.SyncPrototypeStrategyMerge, dialogTestInputs.Next())
			return err
		},
	}
}
