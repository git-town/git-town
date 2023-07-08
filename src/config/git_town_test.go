package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestGitTown(t *testing.T) {
	t.Parallel()

	t.Run("OriginURL()", func(t *testing.T) {
		t.Parallel()
		tests := map[string]giturl.Parts{
			"http://github.com/organization/repository":                     {Host: "github.com", Org: "organization", Repo: "repository"},
			"http://github.com/organization/repository.git":                 {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://github.com/organization/repository":                    {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://github.com/organization/repository.git":                {Host: "github.com", Org: "organization", Repo: "repository"},
			"https://sub.domain.customhost.com/organization/repository":     {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository"},
			"https://sub.domain.customhost.com/organization/repository.git": {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository"},
			"https://dev.azure.com/organization/repository":                 {Host: "dev.azure.com", Org: "organization", Repo: "repository"},
		}
		for give, want := range tests {
			repo := testruntime.CreateGitTown(t)
			os.Setenv("GIT_TOWN_REMOTE", give)
			defer os.Unsetenv("GIT_TOWN_REMOTE")
			have := repo.Config.OriginURL()
			assert.Equal(t, want, *have, give)
		}
	})

	t.Run(".SetOffline()", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
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
