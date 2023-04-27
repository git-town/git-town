package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestRemoveBranch(t *testing.T) {
	t.Parallel()
	runtime := runtime.Create(t)
	err := commands.CreateBranch(&runtime, "b1", "initial")
	assert.NoError(t, err)
	branches, err := runtime.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial", "b1"}, branches)
	err = commands.RemoveBranch(&runtime, "b1")
	assert.NoError(t, err)
	branches, err = runtime.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial"}, branches)
}
