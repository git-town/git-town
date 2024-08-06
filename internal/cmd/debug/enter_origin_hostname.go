package debug

import (
	"os"

	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterOriginHostname() *cobra.Command {
	return &cobra.Command{
		Use: "origin-hostname",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.OriginHostname(configdomain.ParseHostingOriginHostname(""), dialogInputs.Next())
			return err
		},
	}
}
