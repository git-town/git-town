package commands_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/commands"
	"github.com/git-town/git-town/v8/test/runtime"
	"github.com/stretchr/testify/assert"
)

func TestHasGitTownConfig(t *testing.T) {
	t.Parallel()
	runtime := runtime.Create(t)
	res := commands.HasGitTownConfigNow(&runtime)
	assert.False(t, res)
	err := commands.CreateBranch(&runtime, "main", "initial")
	assert.NoError(t, err)
	err = runtime.CreateFeatureBranch("foo")
	assert.NoError(t, err)
	res = commands.HasGitTownConfigNow(&runtime)
	assert.NoError(t, err)
	assert.True(t, res)
}
