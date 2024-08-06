package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestContainsAny(t *testing.T) {
	t.Parallel()

	t.Run("no elements in common", func(t *testing.T) {
		t.Parallel()
		haystack := []int{1, 2, 3}
		needles := []int{4, 5, 6}
		must.False(t, slice.ContainsAny(haystack, needles))
	})

	t.Run("one element in common", func(t *testing.T) {
		t.Parallel()
		haystack := []int{1, 2, 3}
		needles := []int{3, 4, 5}
		must.True(t, slice.ContainsAny(haystack, needles))
	})

	t.Run("multiple elements in common", func(t *testing.T) {
		t.Parallel()
		haystack := []int{1, 2, 3}
		needles := []int{2, 3, 4}
		must.True(t, slice.ContainsAny(haystack, needles))
	})
}
