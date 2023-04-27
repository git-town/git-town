package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestRemoveBranch(t *testing.T) {
	t.Parallel()
	runtime := repo.Create(t)
	err := repo.CreateBranch(&runtime, "b1", "initial")
	assert.NoError(t, err)
	branches, err := runtime.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial", "b1"}, branches)
	err = repo.RemoveBranch(&runtime, "b1")
	assert.NoError(t, err)
	branches, err = runtime.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial"}, branches)
}
