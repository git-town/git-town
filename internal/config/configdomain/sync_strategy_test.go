package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestNewSyncStrategy(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.SyncStrategy]{
			"":       None[configdomain.SyncStrategy](),
			" ":      None[configdomain.SyncStrategy](),
			"merge":  Some(configdomain.SyncStrategyMerge),
			"merge ": Some(configdomain.SyncStrategyMerge),
			" merge": Some(configdomain.SyncStrategyMerge),
			"Merge":  Some(configdomain.SyncStrategyMerge),
			"MERGE":  Some(configdomain.SyncStrategyMerge),
			"rebase": Some(configdomain.SyncStrategyRebase),
			"Rebase": Some(configdomain.SyncStrategyRebase),
			"REBASE": Some(configdomain.SyncStrategyRebase),
		}
		for give, want := range tests {
			have, err := configdomain.ParseSyncStrategy(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.ParseSyncStrategy("zonk")
		must.Error(t, err)
	})
}
