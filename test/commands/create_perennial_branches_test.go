package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestCreatePerennialBranches(t *testing.T) {
	t.Parallel()
	runtime := runtime.CreateGitTown(t)
	err := commands.CreatePerennialBranches(&runtime, "p1", "p2")
	assert.NoError(t, err)
	branches, err := runtime.LocalBranchesMainFirst("main")
	assert.NoError(t, err)
	assert.Equal(t, []string{"main", "initial", "p1", "p2"}, branches)
	runtime.Reload()
	assert.True(t, runtime.IsPerennialBranch("p1"))
	assert.True(t, runtime.IsPerennialBranch("p2"))
}
