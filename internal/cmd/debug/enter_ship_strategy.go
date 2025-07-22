package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterShipStrategy() *cobra.Command {
	return &cobra.Command{
		Use: "ship-strategy",
		RunE: func(_ *cobra.Command, _ []string) error {
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.ShipStrategy(dialog.Args[configdomain.ShipStrategy]{
				Global: None[configdomain.ShipStrategy](),
				Inputs: dialogTestInputs,
				Local:  None[configdomain.ShipStrategy](),
			})
			return err
		},
	}
}
