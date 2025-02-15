package debug

import (
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog"
	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterHostingPlatform() *cobra.Command {
	return &cobra.Command{
		Use: "hosting-platform",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.HostingPlatform(None[configdomain.ForgeType](), dialogInputs.Next())
			return err
		},
	}
}
