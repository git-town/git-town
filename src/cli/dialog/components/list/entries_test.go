package list_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestEntries(t *testing.T) {
	t.Parallel()

	t.Run("AllDisabled", func(t *testing.T) {
		t.Parallel()
		t.Run("all entries are disabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false}, //nolint:exhaustruct
				{Enabled: false}, //nolint:exhaustruct
			}
			must.True(t, entries.AllDisabled())
		})
		t.Run("some entries are enabled", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false}, //nolint:exhaustruct
				{Enabled: false}, //nolint:exhaustruct
				{Enabled: true},  //nolint:exhaustruct
			}
			must.False(t, entries.AllDisabled())
		})
	})

	t.Run("IndexWithText", func(t *testing.T) {
		t.Parallel()
		t.Run("element is in collection", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three")
			found, have := entries.IndexWithText("two")
			must.True(t, found)
			must.EqOp(t, 1, have)
		})
		t.Run("element is not in collection", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two")
			found, _ := entries.IndexWithText("zonk")
			must.False(t, found)
		})
	})

	t.Run("IndexWithTextOr", func(t *testing.T) {
		t.Parallel()
		t.Run("element is in collection", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three")
			have := entries.IndexWithTextOr("two", 0)
			must.EqOp(t, 1, have)
		})
		t.Run("element is not in collection", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two")
			have := entries.IndexWithTextOr("zonk", 5)
			must.EqOp(t, 5, have)
		})
	})
}
