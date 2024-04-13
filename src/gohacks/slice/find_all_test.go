package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestFindAll(t *testing.T) {
	t.Parallel()

	t.Run("list contains the element", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 1, 3, 1}
		have := slice.FindAll(list, 1)
		want := []int{0, 2, 4}
		must.Eq(t, want, have)
	})

	t.Run("list does not contain the element", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3}
		have := slice.FindAll(list, 4)
		want := []int{}
		must.Eq(t, want, have)
	})
}
