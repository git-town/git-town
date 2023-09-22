package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/undo"
	"github.com/stretchr/testify/assert"
)

func TestStashSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		before := undo.StashSnapshot{
			Amount: 1,
		}
		after := undo.StashSnapshot{
			Amount: 3,
		}
		have := before.Diff(after)
		want := undo.StashDiff{
			EntriesAdded: 2,
		}
		assert.Equal(t, want, have)
	})
}
