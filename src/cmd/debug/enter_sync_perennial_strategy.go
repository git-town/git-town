package debug

import (
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncPerennialStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-perennial-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncPerennialStrategy(configdomain.SyncPerennialStrategyRebase, dialogTestInputs.Next())
			return err
		},
	}
}
