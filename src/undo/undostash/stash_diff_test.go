package undostash_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/undo/undostash"
	"github.com/shoenig/test/must"
)

func TestStashDiff(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("entries added", func(t *testing.T) {
			t.Parallel()
			before := gitdomain.StashSize(1)
			after := gitdomain.StashSize(3)
			have := undostash.NewStashDiff(before, after)
			want := undostash.StashDiff{
				EntriesAdded: 2,
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
