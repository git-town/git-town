package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestNewSyncPerennialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]configdomain.SyncPerennialStrategy{
			"merge":  configdomain.SyncPerennialStrategyMerge,
			"rebase": configdomain.SyncPerennialStrategyRebase,
		}
		for give, want := range tests {
			have, err := configdomain.NewSyncPerennialStrategy(give)
			must.NoError(t, err)
			must.EqOp(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"merge", "Merge", "MERGE"} {
			have, err := configdomain.NewSyncPerennialStrategy(give)
			must.NoError(t, err)
			must.EqOp(t, configdomain.SyncPerennialStrategyMerge, have)
		}
	})

	t.Run("defaults to rebase", func(t *testing.T) {
		t.Parallel()
		have, err := configdomain.NewSyncPerennialStrategy("")
		must.NoError(t, err)
		must.EqOp(t, configdomain.SyncPerennialStrategyRebase, have)
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.NewSyncPerennialStrategy("zonk")
		must.Error(t, err)
	})
}
