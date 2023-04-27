package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestCreateBranch(t *testing.T) {
	t.Parallel()
	t.Run("simple branch name", func(t *testing.T) {
		t.Parallel()
		runtime := repo.Create(t)
		err := repo.CreateBranch(&runtime, "branch1", "initial")
		assert.NoError(t, err)
		currentBranch, err := runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", currentBranch)
		branches, err := runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "branch1"}, branches)
	})

	t.Run("branch name with slashes", func(t *testing.T) {
		t.Parallel()
		runtime := repo.Create(t)
		err := repo.CreateBranch(&runtime, "my/feature", "initial")
		assert.NoError(t, err)
		currentBranch, err := runtime.CurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "initial", currentBranch)
		branches, err := runtime.LocalBranchesMainFirst("initial")
		assert.NoError(t, err)
		assert.Equal(t, []string{"initial", "my/feature"}, branches)
	})
}
