package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/spf13/cobra"
)

func enterShipDeleteTrackingBranch() *cobra.Command {
	return &cobra.Command{
		Use: "ship-delete-tracking-branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialogs.ShipDeleteTrackingBranch(true, dialogTestInputs.Next())
			return err
		},
	}
}
