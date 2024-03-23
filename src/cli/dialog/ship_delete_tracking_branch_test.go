package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("ShipDeleteTrackingBranchEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialog.ShipDeleteTrackingBranchEntryYes.Short())
			must.Eq(t, "no", dialog.ShipDeleteTrackingBranchEntryNo.Short())
		})
	})
}
