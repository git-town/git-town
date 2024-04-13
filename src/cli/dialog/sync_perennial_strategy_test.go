package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestSyncPerennialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SyncPerennialStrategyEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "merge", dialog.SyncPerennialStrategyEntryMerge.Short())
			must.Eq(t, "rebase", dialog.SyncPerennialStrategyEntryRebase.Short())
		})
	})
}
