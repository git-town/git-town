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
		t.Run("more than a page before the end of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 1)
			l.MovePageDown()
			must.EqOp(t, 11, l.Cursor)
		})
		t.Run("less than a page before the end of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, 9)
			l.MovePageDown()
			must.EqOp(t, end, l.Cursor)
		})
		t.Run("at end of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, end)
			l.MovePageDown()
			must.EqOp(t, end, l.Cursor)
		})
		t.Run("first and second entry below are disabled, the next ones are enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 0)
			l.MovePageDown()
			must.EqOp(t, 10, l.Cursor)
		})
		t.Run("all entries below are disabled except the next one", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, 0)
			l.MovePageDown()
			must.EqOp(t, end, l.Cursor)
		})
		t.Run("all entries below are disabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 1)
			l.MovePageDown()
			must.EqOp(t, 1, l.Cursor)
		})
		t.Run("only one enabled entry in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MovePageDown()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("no enabled entries in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MovePageDown()
			must.EqOp(t, 0, l.Cursor)
		})
	})

	t.Run("MovePageUp", func(t *testing.T) {
		t.Parallel()
		t.Run("more than a page before the start of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 12)
			l.MovePageUp()
			must.EqOp(t, 2, l.Cursor)
		})
		t.Run("less than a page before the start of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 3)
			l.MovePageUp()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("at the start of the list", func(t *testing.T) {
			t.Parallel()
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
			}
			l := list.NewList(entries, 0)
			l.MovePageUp()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("first and second entry above are disabled, the next ones are enabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			l := list.NewList(entries, 12)
			l.MovePageUp()
			must.EqOp(t, 2, l.Cursor)
		})
		t.Run("all entries above are disabled except the first one", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: true},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, end)
			l.MovePageUp()
			must.EqOp(t, 0, l.Cursor)
		})
		t.Run("all entries above are disabled", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, end)
			l.MovePageUp()
			must.EqOp(t, end-1, l.Cursor)
		})
		t.Run("only one enabled entry in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
				{Enabled: true},
			}
			end := len(entries) - 1
			l := list.NewList(entries, end)
			l.MovePageUp()
			must.EqOp(t, end, l.Cursor)
		})
		t.Run("no enabled entries in list", func(t *testing.T) {
			entries := list.Entries[configdomain.HostingOriginHostname]{
				{Enabled: false},
				{Enabled: false},
				{Enabled: false},
			}
			l := list.NewList(entries, 0)
			l.MovePageUp()
			must.EqOp(t, 0, l.Cursor)
		})
	})
}
