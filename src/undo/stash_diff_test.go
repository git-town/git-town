package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/undo"
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
			have := undo.NewStashDiff(before, after)
			want := undo.StashDiff{
				EntriesAdded: 2,
			}
			must.EqOp(t, want, have)
		})
		t.Run("no entries added", func(t *testing.T) {
			t.Parallel()
			before := domain.StashSnapshot(1)
			after := domain.StashSnapshot(1)
			have := undo.NewStashDiff(before, after)
			want := undo.StashDiff{
				EntriesAdded: 0,
			}
			must.EqOp(t, want, have)
		})
	})
}
