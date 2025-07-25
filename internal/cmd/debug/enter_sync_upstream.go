package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterSyncUpstream() *cobra.Command {
	return &cobra.Command{
		Use: "sync-upstream",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.SyncUpstream(dialog.Args[configdomain.SyncUpstream]{
				Global: None[configdomain.SyncUpstream](),
				Inputs: inputs,
				Local:  None[configdomain.SyncUpstream](),
			})
			return err
		},
	}
}
