package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/git/giturl"
	"github.com/git-town/git-town/v12/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestGitTown(t *testing.T) {
	t.Parallel()

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		repo.CreateFeatureBranch(gitdomain.NewLocalBranchName("feature1"))
		repo.CreateFeatureBranch(gitdomain.NewLocalBranchName("feature2"))
		repo.Reload()
		have := repo.Lineage
		want := configdomain.Lineage{
			gitdomain.NewLocalBranchName("feature1"): gitdomain.NewLocalBranchName("main"),
			gitdomain.NewLocalBranchName("feature2"): gitdomain.NewLocalBranchName("main"),
		}
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
			have := repo.OriginURL()
			must.EqOp(t, want, *have)
		}
	})

	t.Run("Reload", func(t *testing.T) {
		t.Parallel()
		t.Run("lineage changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			branch := gitdomain.NewLocalBranchName("branch-1")
			repo.CreateFeatureBranch(branch)
			repo.Reload()
			want := configdomain.Lineage{
				branch: gitdomain.NewLocalBranchName("main"),
			}
			must.Eq(t, want, repo.Lineage)
		})
	})

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.SetOffline(true)
		must.NoError(t, err)
		offline := repo.Offline
		must.True(t, offline.Bool())
		err = repo.SetOffline(false)
		must.NoError(t, err)
		offline = repo.Offline
		must.False(t, offline.Bool())
	})
}
