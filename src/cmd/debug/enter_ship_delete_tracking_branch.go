package debug

import (
	"os"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/spf13/cobra"
)

func enterShipDeleteTrackingBranch() *cobra.Command {
	return &cobra.Command{
		Use: "ship-delete-tracking-branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialogscreens.EnterShipDeleteTrackingBranch(true, dialogTestInputs.Next())
			return err
		},
	}
}
