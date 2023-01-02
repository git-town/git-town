package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/test"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	t.Parallel()
	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := test.CreateTestGitTownRepo(t)
		err := repo.Config.Offline.Enable(true)
		assert.NoError(t, err)
		offline := repo.Config.Offline.Enabled()
		assert.True(t, offline)
		err = repo.Config.Offline.Enable(false)
		assert.NoError(t, err)
		offline = repo.Config.Offline.Enabled()
		assert.False(t, offline)
	})
}
