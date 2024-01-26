package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/spf13/cobra"
)

func enterSyncPerennialStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-perennial-strategy",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialogscreens.EnterSyncPerennialStrategy(configdomain.SyncPerennialStrategyRebase, dialogTestInputs.Next())
			return err
		},
	}
}
