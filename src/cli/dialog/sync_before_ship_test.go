package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestSyncBeforeShip(t *testing.T) {
	t.Parallel()

	t.Run("SyncBeforeShipEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialog.SyncBeforeShipEntryYes.Short())
			must.Eq(t, "no", dialog.SyncBeforeShipEntryNo.Short())
		})
	})
}
