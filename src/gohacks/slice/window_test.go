package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestWindow(t *testing.T) {
	t.Parallel()
	t.Run("no elements", func(t *testing.T) {
		t.Parallel()
		have, cursorRow := slice.Window(slice.WindowArgs[int]{
			Elements: []int{},
			Cursor:   0,
			Size:     5,
		})
		must.Eq(t, []int{}, have)
		must.EqOp(t, 0, cursorRow)
	})
	t.Run("one element", func(t *testing.T) {
		t.Parallel()
		have, cursorRow := slice.Window(slice.WindowArgs[int]{
			Elements: []int{1},
			Cursor:   0,
			Size:     5,
		})
		must.Eq(t, []int{1}, have)
		must.EqOp(t, 0, cursorRow)
	})

	t.Run("fewer elements than window size", func(t *testing.T) {
		t.Parallel()
		elements := []int{1, 2, 3, 4}
		t.Run("cursor at first element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   0,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 0, cursorRow)
		})
		t.Run("cursor at second element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   1,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 1, cursorRow)
		})
		t.Run("cursor at third element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   2,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 2, cursorRow)
		})
		t.Run("cursor at middle element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   1,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 1, cursorRow)
		})
		t.Run("cursor at second to last element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   5,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 5, cursorRow)
		})
		t.Run("cursor at last element", func(t *testing.T) {
			t.Parallel()
			have, cursorRow := slice.Window(slice.WindowArgs[int]{
				Elements: elements,
				Cursor:   6,
				Size:     7,
			})
			must.Eq(t, elements, have)
			must.EqOp(t, 6, cursorRow)
		})
	})
	t.Run("more elements than window size", func(t *testing.T) {
		t.Parallel()
		t.Run("cursor at first element", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("cursor at second element", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("cursor at third element", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("cursor at middle element", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("cursor at second to last element", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("cursor at last element", func(t *testing.T) {
			t.Parallel()
		})
	})

}
