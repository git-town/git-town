package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

const shipStrategyLong = "strategy"

// type-safe access to the CLI arguments of type configdomain.ShipStrategy
func ShipStrategy() (AddFunc, ReadShipStrategyFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(shipStrategyLong, "s", "", "override the ship-strategy")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.ShipStrategy], error) {
		value, err := cmd.Flags().GetString(shipStrategyLong)
		if err != nil {
			return None[configdomain.ShipStrategy](), err
		}
		strategyOpt, err := configdomain.ParseShipStrategy(value)
		if err != nil {
			return None[configdomain.ShipStrategy](), err
		}
		strategy, hasStrategy := strategyOpt.Get()
		if !hasStrategy {
			return None[configdomain.ShipStrategy](), nil
		}
		return Some(strategy), nil
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadShipStrategyFunc func(*cobra.Command) (Option[configdomain.ShipStrategy], error)
