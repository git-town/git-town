package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestGitTown(t *testing.T) {
	t.Parallel()

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		must.NoError(t, repo.CreateFeatureBranch(domain.NewLocalBranchName("feature1")))
		must.NoError(t, repo.CreateFeatureBranch(domain.NewLocalBranchName("feature2")))
		repo.GitTown.Reload()
		have := repo.GitTown.Lineage(repo.GitTown.RemoveLocalConfigValue)
		want := configdomain.Lineage{}
		want[domain.NewLocalBranchName("feature1")] = domain.NewLocalBranchName("main")
		want[domain.NewLocalBranchName("feature2")] = domain.NewLocalBranchName("main")
		must.Eq(t, want, have)
	})

	t.Run("OriginURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]giturl.Parts{
			"http://github.com/organization/repository":                     {Host: "github.com", Org: "organization", Repo: "repository", User: ""},
			"http://github.com/organization/repository.git":                 {Host: "github.com", Org: "organization", Repo: "repository", User: ""},
			"https://github.com/organization/repository":                    {Host: "github.com", Org: "organization", Repo: "repository", User: ""},
			"https://github.com/organization/repository.git":                {Host: "github.com", Org: "organization", Repo: "repository", User: ""},
			"https://sub.domain.customhost.com/organization/repository":     {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository", User: ""},
			"https://sub.domain.customhost.com/organization/repository.git": {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository", User: ""},
		}
		for give, want := range tests {
			repo := testruntime.CreateGitTown(t)
			os.Setenv("GIT_TOWN_REMOTE", give)
			defer os.Unsetenv("GIT_TOWN_REMOTE")
			have := repo.GitTown.OriginURL()
			must.EqOp(t, want, *have)
		}
	})

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.GitTown.SetOffline(true)
		must.NoError(t, err)
		offline, err := repo.GitTown.IsOffline()
		must.NoError(t, err)
		must.True(t, offline.Bool())
		err = repo.GitTown.SetOffline(false)
		must.NoError(t, err)
		offline, err = repo.GitTown.IsOffline()
		must.NoError(t, err)
		must.False(t, offline.Bool())
	})
}
