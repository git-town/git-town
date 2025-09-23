package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestTruncateLast(t *testing.T) {
	t.Parallel()

	t.Run("list contains multiple elements", func(t *testing.T) {
		t.Parallel()
		list := []int{1, 2, 3}
		have := slice.TruncateLast(list)
		want := []int{1, 2}
		must.Eq(t, want, have)
	})

	t.Run("list contains no elements", func(t *testing.T) {
		t.Parallel()
		list := []int{}
		have := slice.TruncateLast(list)
		must.Len(t, 0, have)
	})

	t.Run("list contains one element", func(t *testing.T) {
		t.Parallel()
		list := []int{1}
		have := slice.TruncateLast(list)
		must.Len(t, 0, have)
	})

	t.Run("slice alias type", func(t *testing.T) {
		t.Parallel()
		list := gitdomain.SHAs{"111111", "222222", "333333"}
		have := slice.TruncateLast(list)
		want := gitdomain.SHAs{"111111", "222222"}
		must.Eq(t, want, have)
	})
}
