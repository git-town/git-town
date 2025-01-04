package flags

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	syncStrategyLong  = "strategy"
	syncStrategyShort = ""
)

// type-safe access to the CLI arguments of type configdomain.SyncStrategy
func SyncStrategy() (AddFunc, ReadSyncStrategyFunc) {
	addFlag := func(cmd *cobra.Command) {
		cmd.Flags().StringP(syncStrategyLong, syncStrategyShort, "", "override the sync-strategy for the current branch")
	}
	readFlag := func(cmd *cobra.Command) (Option[configdomain.SyncStrategy], error) {
		value, err := cmd.Flags().GetString(syncStrategyLong)
		if err != nil {
			return None[configdomain.SyncStrategy](), err
		}
		strategyOpt, err := configdomain.ParseSyncStrategy(value)
		if err != nil {
			return None[configdomain.SyncStrategy](), err
		}
		strategy, hasStrategy := strategyOpt.Get()
		if !hasStrategy {
			return None[configdomain.SyncStrategy](), nil
		}
		return Some(strategy), nil
	}
	return addFlag, readFlag
}

// ReadSyncStrategyFunc is the type signature for the function that reads the sync "strategy" flag from the args to the given Cobra command.
type ReadSyncStrategyFunc func(*cobra.Command) (Option[configdomain.SyncStrategy], error)
