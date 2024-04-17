package list_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestList(t *testing.T) {
	t.Parallel()
	t.Run("MoveCursorDown", func(t *testing.T) {
		t.Parallel()
		t.Run("at beginning of list", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three")
			have := list.NewList(entries, 0)
			have.MoveCursorDown()
			must.EqOp(t, 1, have.Cursor)
		})
		t.Run("in middle of list", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three", "four")
			have := list.NewList(entries, 1)
			have.MoveCursorDown()
			must.EqOp(t, 2, have.Cursor)
		})
		t.Run("at end of list", func(t *testing.T) {
			t.Parallel()
			entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three")
			have := list.NewList(entries, 2)
			have.MoveCursorDown()
			must.EqOp(t, 0, have.Cursor)
		})
	})
}
