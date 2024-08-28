package debug

import (
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/spf13/cobra"
)

func enterShipStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "ship-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := components.LoadTestInputs(os.Environ())
			_, _, err := dialog.ShipStrategy(configdomain.ShipStrategyAPI, dialogTestInputs.Next())
			return err
		},
	}
}
