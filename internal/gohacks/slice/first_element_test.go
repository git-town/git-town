package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFirstElement(t *testing.T) {
	t.Parallel()
	t.Run("normal slice", func(t *testing.T) {
		t.Parallel()
		ints := []int{1, 2}
		have := slice.FirstElement(ints)
		must.True(t, have.EqualSome(1))
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		ints := []int{}
		have := slice.FirstElement(ints)
		must.True(t, have.IsNone())
	})
}
