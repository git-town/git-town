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
		elements := []int{}
		have := slice.Window(elements, 0, 4)
		want := []int{}
		must.Eq(t, want, have)
	})
	t.Run("one element", func(t *testing.T) {
		t.Parallel()
	})
	t.Run("fewer elements than window size", func(t *testing.T) {
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
