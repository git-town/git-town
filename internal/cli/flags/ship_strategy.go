package flags

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/spf13/cobra"
)

const shipStrategyLong = "strategy"

// type-safe access to the CLI arguments of type configdomain.ShipIntoNonPerennialParentLong
func ShipStrategy() (AddFunc, ReadShipStrategyFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().BoolP(shipStrategyLong, "s", false, "override the ship-strategy")
	}
	readFlag := func(cmd *cobra.Command) (configdomain.ShipStrategy, error) {
		value, err := cmd.Flags().GetString(shipStrategyLong)
		if err != nil {
			return configdomain.ShipStrategyAPI, err
		}
		strategyOpt, err := configdomain.ParseShipStrategy(value)
		if err != nil {
			return configdomain.ShipStrategyAPI, err
		}
		strategy, hasStrategy := strategyOpt.Get()
		if !hasStrategy {
			return configdomain.ShipStrategyAPI, errors.New(messages.ShipStrategyMissing)
		}
		return strategy, nil
	}
	return addFlag, readFlag
}

// the type signature for the function that reads the dry-run flag from the args to the given Cobra command
type ReadShipStrategyFunc func(*cobra.Command) (configdomain.ShipStrategy, error)
