package repo_test

import (
	"testing"

	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/test/repo"
	"github.com/stretchr/testify/assert"
)

func TestAddRemote(t *testing.T) {
	t.Parallel()
	t.Run("remote doesn't exist", func(t *testing.T) {
		t.Parallel()
		dev := repo.Create(t)
		remotes, err := dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{}, remotes)
		origin := repo.Create(t)
		err = repo.AddRemote(&dev, config.OriginRemote, origin.Dir())
		assert.NoError(t, err)
		remotes, err = dev.Remotes()
		assert.NoError(t, err)
		assert.Equal(t, []string{"origin"}, remotes)
	})
}
