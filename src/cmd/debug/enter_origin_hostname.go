package debug

import (
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterOriginHostname() *cobra.Command {
	return &cobra.Command{
		Use: "origin-hostname",
		RunE: func(_ *cobra.Command, args []string) error {
			dialogInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.OriginHostname("", dialogInputs.Next())
			return err
		},
	}
}
