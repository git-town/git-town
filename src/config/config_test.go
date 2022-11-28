package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t)
		err := repo.Config.SetOffline(true)
		assert.NoError(t, err)
		offline := repo.Config.IsOffline()
		assert.True(t, offline)
		err = repo.Config.SetOffline(false)
		assert.NoError(t, err)
		offline = repo.Config.IsOffline()
		assert.False(t, offline)
	})
}
