package config_test

import (
	"testing"

	"github.com/git-town/git-town/v15/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestUnvalidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.Config.SetOffline(true)
		must.NoError(t, err)
		must.True(t, repo.Config.Config.Offline.IsTrue())
		err = repo.Config.SetOffline(false)
		must.NoError(t, err)
		must.False(t, repo.Config.Config.Offline.IsTrue())
	})
}
