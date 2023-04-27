package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestHasGitTownConfig(t *testing.T) {
	t.Parallel()
	runtime := repo.Create(t)
	res := repo.HasGitTownConfigNow(&runtime)
	assert.False(t, res)
	err := repo.CreateBranch(&runtime, "main", "initial")
	assert.NoError(t, err)
	err = runtime.CreateFeatureBranch("foo")
	assert.NoError(t, err)
	res = repo.HasGitTownConfigNow(&runtime)
	assert.NoError(t, err)
	assert.True(t, res)
}
