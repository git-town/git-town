package git_test

import (
	"testing"

	"github.com/git-town/git-town/src/git"
	"github.com/stretchr/testify/assert"
)

func TestCurrentBranchTracker(t *testing.T) {
	cbc := git.CurrentBranchCache{}
	cbc.Set("foo")
	assert.Equal(t, "foo", cbc.Current())
}
