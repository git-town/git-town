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
		entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three", "four")
		tests := map[int]int{
			0: 1, // at beginning of list
			1: 2, // in middle of list
			3: 0, // at end of list
		}
		for give, want := range tests {
			have := list.NewList(entries, give)
			have.MoveCursorDown()
			must.EqOp(t, want, have.Cursor)
		}
	})

	t.Run("MoveCursorUp", func(t *testing.T) {
		t.Parallel()
		entries := list.NewEntries[configdomain.HostingOriginHostname]("one", "two", "three", "four")
		tests := map[int]int{
			0: 3, // at beginning of list
			2: 1, // in middle of list
			3: 2, // at end of list
		}
		for give, want := range tests {
			have := list.NewList(entries, give)
			have.MoveCursorUp()
			must.EqOp(t, want, have.Cursor)
		}
	})
}
