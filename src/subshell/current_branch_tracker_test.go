package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/stretchr/testify/assert"
)

func TestCurrentBranchTracker(t *testing.T) {
	t.Parallel()
	t.Run("git checkout <branch>", func(t *testing.T) {
		t.Parallel()
		tracker := subshell.CurrentBranchTracker{Value: "first"}
		tracker.Track("git", "checkout", "second")
		assert.Equal(t, tracker.Value, "second")
	})
	t.Run("other Git command", func(t *testing.T) {
		t.Parallel()
		tracker := subshell.CurrentBranchTracker{Value: "first"}
		tracker.Track("git", "status")
		assert.Equal(t, tracker.Value, "first")
	})
}
