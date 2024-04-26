package debug

import (
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/spf13/cobra"
)

func enterHostingPlatform() *cobra.Command {
	return &cobra.Command{
		Use: "hosting-platform",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.HostingPlatform(None[configdomain.HostingPlatform](), dialogInputs.Next())
			return err
		},
	}
}
