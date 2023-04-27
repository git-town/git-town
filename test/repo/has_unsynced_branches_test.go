package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/fixture"
	"github.com/git-town/git-town/v8/test/fs"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestHasUnsyncedBranches(t *testing.T) {
	t.Parallel()
	t.Run("branches are in sync", func(t *testing.T) {
		t.Parallel()
		env, err := fixture.NewStandardFixture(t.TempDir())
		assert.NoError(t, err)
		runner := env.DevRepo
		err = repo.CreateBranch(&runner, "branch1", "main")
		assert.NoError(t, err)
		err = runner.CheckoutBranch("branch1")
		assert.NoError(t, err)
		err = fs.CreateFile(runner.Dir(), "file1", "content")
		assert.NoError(t, err)
		err = repo.StageFiles(&runner, "file1")
		assert.NoError(t, err)
		err = repo.CommitStagedChanges(&runner, "stuff")
		assert.NoError(t, err)
		err = repo.PushBranchToRemote(&runner, "branch1", config.OriginRemote)
		assert.NoError(t, err)
		have, err := repo.HasUnsyncedBranches(&runner)
		assert.NoError(t, err)
		assert.False(t, have)
	})

	t.Run("branch is ahead", func(t *testing.T) {
		t.Parallel()
		env, err := fixture.NewStandardFixture(t.TempDir())
		assert.NoError(t, err)
		err = repo.CreateBranch(&env.DevRepo, "branch1", "main")
		assert.NoError(t, err)
		err = repo.PushBranch(&env.DevRepo)
		assert.NoError(t, err)
		err = fs.CreateFile(env.DevRepo.Dir(), "file1", "content")
		assert.NoError(t, err)
		err = repo.StageFiles(&env.DevRepo, "file1")
		assert.NoError(t, err)
		err = repo.CommitStagedChanges(&env.DevRepo, "stuff")
		assert.NoError(t, err)
		have, err := repo.HasUnsyncedBranches(&env.DevRepo)
		assert.NoError(t, err)
		assert.True(t, have)
	})

	t.Run("branch is behind", func(t *testing.T) {
		t.Parallel()
		env, err := fixture.NewStandardFixture(t.TempDir())
		assert.NoError(t, err)
		err = repo.CreateBranch(&env.DevRepo, "branch1", "main")
		assert.NoError(t, err)
		err = repo.PushBranch(&env.DevRepo)
		assert.NoError(t, err)
		err = env.OriginRepo.CheckoutBranch("main")
		assert.NoError(t, err)
		err = fs.CreateFile(env.OriginRepo.Dir(), "file1", "content")
		assert.NoError(t, err)
		err = repo.StageFiles(env.OriginRepo, "file1")
		assert.NoError(t, err)
		err = repo.CommitStagedChanges(env.OriginRepo, "stuff")
		assert.NoError(t, err)
		err = env.OriginRepo.CheckoutBranch("initial")
		assert.NoError(t, err)
		err = repo.Fetch(&env.DevRepo)
		assert.NoError(t, err)
		have, err := repo.HasUnsyncedBranches(&env.DevRepo)
		assert.NoError(t, err)
		assert.True(t, have)
	})
}
