package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestSyncBeforeShip(t *testing.T) {
	t.Parallel()

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()
		give := configdomain.NewSyncBeforeShip(true)
		have := give.Bool()
		must.True(t, have)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		give := configdomain.NewSyncBeforeShip(true)
		have := give.String()
		want := "true"
		must.EqOp(t, want, have)
	})

	t.Run("NewSyncBeforeShip", func(t *testing.T) {
		t.Parallel()
		have := configdomain.NewSyncBeforeShip(true)
		want := configdomain.SyncBeforeShip(true)
		must.EqOp(t, want, have)
	})

	t.Run("NewSyncBeforeShipRef", func(t *testing.T) {
		t.Parallel()
		have := configdomain.NewSyncBeforeShipRef(true)
		want := configdomain.SyncBeforeShip(true)
		must.EqOp(t, want, *have)
	})

	t.Run("ParseSyncBeforeShip", func(t *testing.T) {
		t.Parallel()
		t.Run("parsable value", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseSyncBeforeShip("yes", "test")
			must.NoError(t, err)
			want := configdomain.NewSyncBeforeShip(true)
			must.EqOp(t, want, have)
		})
		t.Run("invalid value", func(t *testing.T) {
			t.Parallel()
			_, err := configdomain.ParseSyncBeforeShip("zonk", "local config")
			must.EqOp(t, `invalid value for local config: "zonk". Please provide either "yes" or "no"`, err.Error())
		})
	})
}
