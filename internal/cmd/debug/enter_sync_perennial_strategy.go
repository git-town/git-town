package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterSyncPerennialStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "sync-perennial-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.SyncPerennialStrategy(dialog.Args[configdomain.SyncPerennialStrategy]{
				Global: None[configdomain.SyncPerennialStrategy](),
				Inputs: dialogTestInputs,
				Local:  Some(configdomain.SyncPerennialStrategyRebase),
			})
			return err
		},
	}
}
