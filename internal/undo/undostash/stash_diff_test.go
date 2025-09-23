package undostash_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undostash"
	"github.com/shoenig/test/must"
)

func TestStashDiff(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("entries added", func(t *testing.T) {
			t.Parallel()
			before := gitdomain.StashSize(1)
			after := gitdomain.StashSize(2)
			have := undostash.NewStashDiff(before, after)
			want := undostash.StashDiff{
				EntriesAdded: 1,
			}
			must.EqOp(t, want, have)
		})
		t.Run("no entries added", func(t *testing.T) {
			t.Parallel()
			before := gitdomain.StashSize(1)
			after := gitdomain.StashSize(1)
			have := undostash.NewStashDiff(before, after)
			want := undostash.StashDiff{
				EntriesAdded: 0,
			}
			must.EqOp(t, want, have)
		})
	})
}
