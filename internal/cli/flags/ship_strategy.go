package flags

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const shipStrategyLong = "strategy"

// type-safe access to the CLI arguments of type configdomain.ShipIntoNonPerennialParentLong
func ShipStrategy() (AddFunc, ReadShipStrategyFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(shipStrategyLong, "s", "", "override the ship-strategy")
	}
	readFlag := func(cmd *cobra.Command) Option[configdomain.ShipStrategy] {
		value, err := cmd.Flags().GetString(shipStrategyLong)
		if err != nil {
			print.Error(err)
			return None[configdomain.ShipStrategy]()
		}
		strategyOpt, err := configdomain.ParseShipStrategy(value)
		if err != nil {
			print.Error(err)
			return None[configdomain.ShipStrategy]()
		}
		strategy, hasStrategy := strategyOpt.Get()
		if !hasStrategy {
			print.Error(errors.New(messages.ShipStrategyMissing))
			return None[configdomain.ShipStrategy]()
		}
		return Some(strategy)
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadShipStrategyFunc func(*cobra.Command) Option[configdomain.ShipStrategy]
