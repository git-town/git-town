package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/spf13/cobra"
)

func enterShipDeleteTrackingBranch() *cobra.Command {
	return &cobra.Command{
		Use: "ship-delete-tracking-branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.ShipDeleteTrackingBranch(true, dialogTestInputs.Next())
			return err
		},
	}
}
