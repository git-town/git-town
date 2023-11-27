package config_test

import (
	"testing"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/shoenig/test/must"
)

func TestNewSyncPerennialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]config.SyncPerennialStrategy{
			"merge":  config.SyncPerennialStrategyMerge,
			"rebase": config.SyncPerennialStrategyRebase,
		}
		for give, want := range tests {
			have, err := config.NewSyncPerennialStrategy(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"merge", "Merge", "MERGE"} {
			have, err := config.NewSyncPerennialStrategy(give)
			must.NoError(t, err)
			must.EqOp(t, config.SyncPerennialStrategyMerge, have)
		}
	})

	t.Run("defaults to rebase", func(t *testing.T) {
		t.Parallel()
		have, err := config.NewSyncPerennialStrategy("")
		must.NoError(t, err)
		must.EqOp(t, config.SyncPerennialStrategyRebase, have)
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()
		_, err := config.NewSyncPerennialStrategy("zonk")
		must.Error(t, err)
	})
}
