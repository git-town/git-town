package dialogs_test

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestSyncUpstream(t *testing.T) {
	t.Parallel()

	t.Run("SyncUpstreamEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogs.SyncUpstreamEntryYes.Short())
			must.Eq(t, "no", dialogs.SyncUpstreamEntryNo.Short())
		})
	})
}
