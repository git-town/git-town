package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestPushBranchToRemote(t *testing.T) {
	t.Parallel()
	dev := repo.Create(t)
	origin := repo.Create(t)
	err := repo.AddRemote(&dev, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = repo.CreateBranch(&dev, "b1", "initial")
	assert.NoError(t, err)
	err = repo.PushBranchToRemote(&dev, "b1", config.OriginRemote)
	assert.NoError(t, err)
	branches, err := origin.LocalBranchesMainFirst("initial")
	assert.NoError(t, err)
	assert.Equal(t, []string{"initial", "b1"}, branches)
}
