package dialogs_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialogs"
	"github.com/shoenig/test/must"
)

func TestShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("ShipDeleteTrackingBranchEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogs.ShipDeleteTrackingBranchEntryYes.Short())
			must.Eq(t, "no", dialogs.ShipDeleteTrackingBranchEntryNo.Short())
		})
	})
}
