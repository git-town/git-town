package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncPerennialStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-perennial-strategy",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.SyncPerennialStrategy(configdomain.SyncPerennialStrategyRebase, dialogTestInputs.Next())
			return err
		},
	}
}
