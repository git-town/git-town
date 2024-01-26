package enter_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/shoenig/test/must"
)

func TestSyncBeforeShip(t *testing.T) {
	t.Parallel()

	t.Run("SyncBeforeShipEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", enter.SyncBeforeShipEntryYes.Short())
			must.Eq(t, "no", enter.SyncBeforeShipEntryNo.Short())
		})
	})
}
