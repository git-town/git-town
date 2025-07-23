package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterSyncTags() *cobra.Command {
	return &cobra.Command{
		Use: "sync-tags",
		RunE: func(_ *cobra.Command, _ []string) error {
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err := dialog.SyncTags(dialog.Args[configdomain.SyncTags]{
				Global: None[configdomain.SyncTags](),
				Inputs: inputs,
				Local:  None[configdomain.SyncTags](),
			})
			return err
		},
	}
}
