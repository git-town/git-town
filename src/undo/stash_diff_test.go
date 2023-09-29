package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/stretchr/testify/assert"
)

func TestStashDiff(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		before := domain.StashSnapshot{
			Amount: 1,
		}
		after := domain.StashSnapshot{
			Amount: 3,
		}
		have := undo.NewStashDiff(before, after)
		want := undo.StashDiff{
			EntriesAdded: 2,
		}
		assert.Equal(t, want, have)
	})
}
