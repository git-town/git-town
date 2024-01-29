package dialogs_test

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestSyncPerennialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SyncPerennialStrategyEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "merge", dialogs.SyncPerennialStrategyEntryMerge.Short())
			must.Eq(t, "rebase", dialogs.SyncPerennialStrategyEntryRebase.Short())
		})
	})
}
