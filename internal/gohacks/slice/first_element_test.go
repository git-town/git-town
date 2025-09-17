package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFirstElement(t *testing.T) {
	t.Parallel()

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		list := []int{}
		have := slice.FirstElement(list)
		must.True(t, have.IsNone())
	})

	t.Run("list with many elements", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2}
		have := slice.FirstElement(list)
		must.True(t, have.EqualSome(1))
	})

	t.Run("list with one element", func(t *testing.T) {
		t.Parallel()
		list := []int{1}
		have := slice.FirstElement(list)
		must.True(t, have.EqualSome(1))
	})

}
