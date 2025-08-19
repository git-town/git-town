package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("IsTrue", func(t *testing.T) {
		t.Parallel()
		give := configdomain.ShipDeleteTrackingBranch(true)
		have := give.IsTrue()
		must.True(t, have)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := configdomain.ShipDeleteTrackingBranch(true)
		have := give.String()
		want := "true"
		must.EqOp(t, want, have)
	})
}
