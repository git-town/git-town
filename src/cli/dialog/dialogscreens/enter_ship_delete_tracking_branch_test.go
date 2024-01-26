package dialogscreens_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/shoenig/test/must"
)

func TestEnterShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("ShipDeleteTrackingBranchEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogscreens.ShipDeleteTrackingBranchEntryYes.Short())
			must.Eq(t, "no", dialogscreens.ShipDeleteTrackingBranchEntryNo.Short())
		})
	})
}
