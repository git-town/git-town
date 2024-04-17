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
		start := 0
		end := len(entries)
		tests := map[int]int{
			start: 1,     // at start of list
			1:     2,     // in middle of list
			end:   start, // at end of list
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
		start := 0
		end := len(entries) - 1
		tests := map[int]int{
			start: end, // at beginning of list
			2:     1,   // in middle of list
			end:   2,   // at end of list
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
		start := 0
		end := len(entries) - 1
		tests := map[int]int{
			start: 10,  // at beginning of list
			10:    end, // less than a page before the end of the list
			end:   end, // at end of list
		}
		for give, want := range tests {
			have := list.NewList(entries, give)
			have.MovePageDown()
			must.EqOp(t, want, have.Cursor)
		}
	})

	t.Run("MovePageUp", func(t *testing.T) {
		t.Parallel()
		entries := list.NewEntries[configdomain.HostingOriginHostname]("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12")
		start := 0
		end := len(entries) - 1
		tests := map[int]int{
			start: start, // at beginning of list
			2:     start, // less than a page before the start of the list
			11:    1,     // more than a page before the start of the list
			end:   2,     // at end of list
		}
		for give, want := range tests {
			have := list.NewList(entries, give)
			have.MovePageUp()
			must.EqOp(t, want, have.Cursor)
		}
	})
}
