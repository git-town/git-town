package dialogs_test

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestSyncBeforeShip(t *testing.T) {
	t.Parallel()

	t.Run("SyncBeforeShipEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogs.SyncBeforeShipEntryYes.Short())
			must.Eq(t, "no", dialogs.SyncBeforeShipEntryNo.Short())
		})
	})
}
