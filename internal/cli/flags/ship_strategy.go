package flags

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
		strategy, hasStrategy := strategyOpt.Get()
		if err != nil || !hasStrategy {
			return None[configdomain.ShipStrategy](), err
		}
		return Some(strategy), nil
	}
	return addFlag, readFlag
}

// ReadShipStrategyFunc is the type signature for the function that reads the ship "strategy" flag from the args to the given Cobra command.
type ReadShipStrategyFunc func(*cobra.Command) (Option[configdomain.ShipStrategy], error)
