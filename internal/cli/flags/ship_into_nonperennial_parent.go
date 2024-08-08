package flags

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/spf13/cobra"
)

const shipIntoNonPerennialParentLong = "to-parent"

// type-safe access to the CLI arguments of type configdomain.ShipIntoNonPerennialParentLong
func ShipIntoNonPerennialParent() (AddFunc, ReadShipIntoNonPerennialParentFlagFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(shipIntoNonPerennialParentLong, "p", false, "allow shipping into non-perennial parent")
	}
	readFlag := func(cmd *cobra.Command) configdomain.ShipIntoNonperennialParent {
		value, err := cmd.Flags().GetBool(shipIntoNonPerennialParentLong)
		if err != nil {
			panic(fmt.Sprintf(messages.FlagStringDoesntExist, cmd.Name(), shipIntoNonPerennialParentLong))
		}
		return configdomain.ShipIntoNonperennialParent(value)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadShipIntoNonPerennialParentFlagFunc func(*cobra.Command) configdomain.ShipIntoNonperennialParent
