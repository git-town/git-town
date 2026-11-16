package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestWindow(t *testing.T) {
	t.Parallel()
	tests := map[slice.WindowArgs]slice.WindowResult{
		// no elements
		{CursorPos: 0, ElementCount: 0, WindowSize: 5}: {EndRow: 0, StartRow: 0},
		// one element
		{CursorPos: 0, ElementCount: 1, WindowSize: 5}: {EndRow: 1, StartRow: 0},

		// FEWER ELEMENTS THAN WINDOW SIZE

		// cursor at first element
		{CursorPos: 0, ElementCount: 7, WindowSize: 9}: {EndRow: 7, StartRow: 0},
		// fewer elements than window size, cursor at second element
		{CursorPos: 1, ElementCount: 7, WindowSize: 9}: {EndRow: 7, StartRow: 0},
		// fewer elements than window size, cursor at middle element
		{CursorPos: 3, ElementCount: 7, WindowSize: 9}: {EndRow: 7, StartRow: 0},
		// fewer elements than window size, cursor at second to last element
		{CursorPos: 5, ElementCount: 7, WindowSize: 9}: {EndRow: 7, StartRow: 0},
		// fewer elements than window size, cursor at last element
		{CursorPos: 6, ElementCount: 7, WindowSize: 9}: {EndRow: 7, StartRow: 0},

		// MORE ELEMENTS THAN WINDOW SIZE

		// cursor at first element
		{CursorPos: 0, ElementCount: 20, WindowSize: 9}: {EndRow: 9, StartRow: 0},
		// fewer elements than window size, cursor at second element
		{CursorPos: 1, ElementCount: 20, WindowSize: 9}: {EndRow: 9, StartRow: 0},
		// fewer elements than window size, cursor at middle element
		{CursorPos: 10, ElementCount: 20, WindowSize: 9}: {EndRow: 15, StartRow: 6},
		// fewer elements than window size, cursor at second to last element
		{CursorPos: 18, ElementCount: 20, WindowSize: 9}: {EndRow: 20, StartRow: 11},
		// fewer elements than window size, cursor at last element
		{CursorPos: 19, ElementCount: 20, WindowSize: 9}: {EndRow: 20, StartRow: 11},
	}
	for give, want := range tests {
		have := slice.Window(give)
		must.EqOp(t, want, have)
	}
}
