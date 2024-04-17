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

	t.Run("MovePageDown", func(t *testing.T) {
		t.Parallel()
		entries := list.NewEntries[configdomain.HostingOriginHostname]("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12")
		tests := map[int]int{
			0:  10, // at beginning of list
			10: 12, // less than a page before the end of the list
			11: 12, // at end of list
		}
		for give, want := range tests {
			have := list.NewList(entries, give)
			have.MovePageDown()
			must.EqOp(t, want, have.Cursor)
		}
	})
}
