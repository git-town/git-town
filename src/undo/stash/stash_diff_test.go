package stash_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo/stash"
	"github.com/shoenig/test/must"
)

func TestStashDiff(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("entries added", func(t *testing.T) {
			t.Parallel()
			before := domain.StashSnapshot(1)
			after := domain.StashSnapshot(3)
			have := stash.NewDiff(before, after)
			want := stash.Diff{
				EntriesAdded: 2,
			}
			must.EqOp(t, want, have)
		})
		t.Run("no entries added", func(t *testing.T) {
			t.Parallel()
			before := domain.StashSnapshot(1)
			after := domain.StashSnapshot(1)
			have := stash.NewDiff(before, after)
			want := stash.Diff{
				EntriesAdded: 0,
			}
			must.EqOp(t, want, have)
		})
	})
}
