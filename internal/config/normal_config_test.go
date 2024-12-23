package config_test

import (
	"testing"

	"github.com/git-town/git-town/v17/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestNormalConfig(t *testing.T) {
	t.Parallel()

	t.Run("RemoteURL", func(t *testing.T) {
		t.Parallel()

	})

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.Config.NormalConfig.SetOffline(true)
		must.NoError(t, err)
		must.True(t, repo.Config.NormalConfig.Offline.IsTrue())
		err = repo.Config.NormalConfig.SetOffline(false)
		must.NoError(t, err)
		must.False(t, repo.Config.NormalConfig.Offline.IsTrue())
	})
}
