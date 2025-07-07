package config_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestNormalConfig(t *testing.T) {
	t.Parallel()

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.Config.NormalConfig.SetOffline(repo.TestRunner, true)
		must.NoError(t, err)
		must.True(t, repo.Config.NormalConfig.Offline.IsOffline())
		err = repo.Config.NormalConfig.SetOffline(repo.TestRunner, false)
		must.NoError(t, err)
		must.False(t, repo.Config.NormalConfig.Offline.IsOffline())
	})
}
