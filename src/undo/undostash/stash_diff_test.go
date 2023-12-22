package undostash_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/undo/undostash"
	"github.com/shoenig/test/must"
)

func TestStashDiff(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("entries added", func(t *testing.T) {
			t.Parallel()
			before := undodomain.StashSnapshot(1)
			after := undodomain.StashSnapshot(3)
			have := undostash.NewStashDiff(before, after)
			want := undostash.StashDiff{
				EntriesAdded: 2,
			}
			must.EqOp(t, want, have)
		})
		t.Run("no entries added", func(t *testing.T) {
			t.Parallel()
			before := undodomain.StashSnapshot(1)
			after := undodomain.StashSnapshot(1)
			have := undostash.NewStashDiff(before, after)
			want := undostash.StashDiff{
				EntriesAdded: 0,
			}
			must.EqOp(t, want, have)
		})
	})
}
