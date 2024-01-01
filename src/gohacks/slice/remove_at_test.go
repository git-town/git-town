package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestRemoveAt(t *testing.T) {
	t.Parallel()

	t.Run("index is within the list", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3}
		slice.RemoveAt(&list, 1)
		want := []int{1, 3}
		must.Eq(t, want, list)
	})

	t.Run("index is at end of list", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3}
		slice.RemoveAt(&list, 2)
		want := []int{1, 2}
		must.Eq(t, want, list)
	})

	t.Run("index is at beginning of list", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3}
		slice.RemoveAt(&list, 0)
		want := []int{2, 3}
		must.Eq(t, want, list)
	})

	t.Run("slice alias type", func(t *testing.T) {
		t.Parallel()
		list := gitdomain.SHAs{gitdomain.NewSHA("111111"), gitdomain.NewSHA("222222"), gitdomain.NewSHA("333333")}
		slice.RemoveAt(&list, 0)
		want := gitdomain.SHAs{gitdomain.NewSHA("222222"), gitdomain.NewSHA("333333")}
		must.Eq(t, want, list)
	})
}
