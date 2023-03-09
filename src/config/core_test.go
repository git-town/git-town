package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestGitTown(t *testing.T) {
	t.Parallel()
	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t)
		err := repo.Config.SetOffline(true)
		assert.NoError(t, err)
		offline, err := repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.True(t, offline)
		err = repo.Config.SetOffline(false)
		assert.NoError(t, err)
		offline, err = repo.Config.IsOffline()
		assert.Nil(t, err)
		assert.False(t, offline)
	})
}
