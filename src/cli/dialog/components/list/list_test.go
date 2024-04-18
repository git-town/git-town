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
		t.Run("entry below is enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorDown()
			must.EqOp(t, 1, l.Cursor)
		})
		t.Run("at end of list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 2)
			l.MoveCursorDown()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("first and second entry below are disabled, the next one is enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorDown()
			must.EqOp(t, 3, l.Cursor)
		})
		t.Run("all entries below are disabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 1)
			l.MoveCursorDown()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("only one enabled entry in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorDown()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("no enabled entries in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorDown()
			must.EqOp(t, 0, l.Cursor)
		})
	})

	t.Run("MoveCursorUp", func(t *testing.T) {
		t.Parallel()
		t.Run("entry above is enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 2)
			l.MoveCursorUp()
			must.EqOp(t, 1, l.Cursor)
		})
		t.Run("at beginning of list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorUp()
			must.EqOp(t, 2, l.Cursor)
		})
		t.Run("first and second entry above are disabled, the next one is enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			l := list.NewList(entries, 3)
			l.MoveCursorUp()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("all entries above are disabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 2)
			l.MoveCursorUp()
			must.EqOp(t, 3, l.Cursor)
		})
		t.Run("only one enabled entry in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorUp()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("no enabled entries in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MoveCursorUp()
			must.EqOp(t, 0, l.Cursor)
		})
	})

	t.Run("MovePageDown", func(t *testing.T) {
		t.Parallel()
		entries := list.NewEntries[configdomain.HostingOriginHostname]("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12")
		start := 0
		end := len(entries) - 1
		tests := map[int]int{
			start: 10,  // at beginning of list
			1:     11,  // more than a page before the end of the list
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
