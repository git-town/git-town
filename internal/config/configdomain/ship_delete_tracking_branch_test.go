package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := configdomain.ShipDeleteTrackingBranch(true)
		have := give.String()
		want := "true"
		must.EqOp(t, want, have)
	})
}
