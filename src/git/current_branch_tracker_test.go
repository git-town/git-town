package git_test

import (
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestCurrentBranchTracker(t *testing.T) {
	cbt := git.CurrentBranchTracker{}
	cbt.Changed("foo")
	assert.Equal(t, "foo", cbt.Current())
}
