package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestShipDeleteTrackingBranch(t *testing.T) {
	t.Parallel()

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()
		give := configdomain.ShipDeleteTrackingBranch(true)
		have := give.Bool()
		must.True(t, have)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := configdomain.ShipDeleteTrackingBranch(true)
		have := give.String()
		want := "true"
		must.EqOp(t, want, have)
	})

	t.Run("ParseShipDeleteTrackingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("parsable value", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseShipDeleteTrackingBranch("yes", "test")
			must.NoError(t, err)
			want := configdomain.ShipDeleteTrackingBranch(true)
			must.EqOp(t, want, have)
		})
		t.Run("invalid value", func(t *testing.T) {
			t.Parallel()
			_, err := configdomain.ParseShipDeleteTrackingBranch("zonk", "local config")
			must.EqOp(t, `invalid value for local config: "zonk". Please provide either "yes" or "no"`, err.Error())
		})
	})
}
