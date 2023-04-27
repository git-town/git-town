package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	t.Parallel()
	dev := repo.Create(t)
	origin := repo.Create(t)
	err := repo.AddRemote(&dev, config.OriginRemote, origin.Dir())
	assert.NoError(t, err)
	err = repo.Fetch(&dev)
	assert.NoError(t, err)
}
