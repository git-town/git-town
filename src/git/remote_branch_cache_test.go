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
	remoteBranches, err := rbt.Get()
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, remoteBranches)
}
