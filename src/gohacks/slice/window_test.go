package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestWindow(t *testing.T) {
	t.Parallel()
	tests := map[slice.WindowArgs]slice.WindowResult{
		// no elements
		{ElementCount: 0, CursorPos: 0, WindowSize: 5}: {StartRow: 0, EndRow: 0},
		// one element
		{ElementCount: 1, CursorPos: 0, WindowSize: 5}: {StartRow: 0, EndRow: 1},

		// FEWER ELEMENTS THAN WINDOW SIZE

		// cursor at first element
		{ElementCount: 7, CursorPos: 0, WindowSize: 9}: {StartRow: 0, EndRow: 7},
		// fewer elements than window size, cursor at second element
		{ElementCount: 7, CursorPos: 1, WindowSize: 9}: {StartRow: 0, EndRow: 7},
		// fewer elements than window size, cursor at middle element
		{ElementCount: 7, CursorPos: 3, WindowSize: 9}: {StartRow: 0, EndRow: 7},
		// fewer elements than window size, cursor at second to last element
		{ElementCount: 7, CursorPos: 5, WindowSize: 9}: {StartRow: 0, EndRow: 7},
		// fewer elements than window size, cursor at last element
		{ElementCount: 7, CursorPos: 6, WindowSize: 9}: {StartRow: 0, EndRow: 7},

		// MORE ELEMENTS THAN WINDOW SIZE

		// cursor at first element
		{ElementCount: 20, CursorPos: 0, WindowSize: 9}: {StartRow: 0, EndRow: 9},
		// fewer elements than window size, cursor at second element
		{ElementCount: 20, CursorPos: 1, WindowSize: 9}: {StartRow: 0, EndRow: 9},
		// fewer elements than window size, cursor at middle element
		{ElementCount: 20, CursorPos: 10, WindowSize: 9}: {StartRow: 6, EndRow: 15},
		// fewer elements than window size, cursor at second to last element
		{ElementCount: 20, CursorPos: 18, WindowSize: 9}: {StartRow: 11, EndRow: 20},
		// fewer elements than window size, cursor at last element
		{ElementCount: 20, CursorPos: 19, WindowSize: 9}: {StartRow: 11, EndRow: 20},
	}
	for give, want := range tests {
		have := slice.Window(give)
		must.EqOp(t, want, have)
	}
}
