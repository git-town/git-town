package git_test

import (
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestRemoteBranchTracker(t *testing.T) {
	rbt := git.RemoteBranchCache{}
	assert.False(t, rbt.Initialized())
	rbt.Set([]string{"one", "two"})
	assert.True(t, rbt.Initialized())
	assert.Equal(t, []string{"one", "two"}, rbt.Get())
}
