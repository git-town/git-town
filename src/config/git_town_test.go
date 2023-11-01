package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git/giturl"
	"github.com/git-town/git-town/v10/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestGitTown(t *testing.T) {
	t.Parallel()

	t.Run("DetermineOriginURL", func(t *testing.T) {
		t.Parallel()
		t.Run("SSH URL", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@github.com:git-town/docs.git", "", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "github.com",
				Org:  "git-town",
				Repo: "docs",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("HTTPS URL", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("https://github.com/git-town/docs.git", "", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "github.com",
				Org:  "git-town",
				Repo: "docs",
				User: "",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("GitLab handbook repo on gitlab.com", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@gitlab.com:gitlab-com/www-gitlab-com.git", "", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "gitlab.com",
				Org:  "gitlab-com",
				Repo: "www-gitlab-com",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("GitLab repository nested inside a group", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@gitlab.com:gitlab-org/quality/triage-ops.git", "", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "gitlab.com",
				Org:  "gitlab-org/quality",
				Repo: "triage-ops",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("self-hosted GitLab server without URL override", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@self-hosted-gitlab.com:git-town/git-town.git", "", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "self-hosted-gitlab.com",
				Org:  "git-town",
				Repo: "git-town",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("self-hosted GitLab server with URL override", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@self-hosted-gitlab.com:git-town/git-town.git", "override.com", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "override.com",
				Org:  "git-town",
				Repo: "git-town",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
		t.Run("custom SSH identity with hostname override", func(t *testing.T) {
			t.Parallel()
			have := config.DetermineOriginURL("git@my-ssh-identity.com:git-town/git-town.git", "gitlab.com", config.OriginURLCache{})
			want := &giturl.Parts{
				Host: "gitlab.com",
				Org:  "git-town",
				Repo: "git-town",
				User: "git",
			}
			must.EqOp(t, *want, *have)
		})
	})

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		must.NoError(t, repo.CreateFeatureBranch(domain.NewLocalBranchName("feature1")))
		must.NoError(t, repo.CreateFeatureBranch(domain.NewLocalBranchName("feature2")))
		repo.Config.Reload()
		have := repo.Config.Lineage(repo.Config.RemoveLocalConfigValue)
		want := config.Lineage{}
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
			have := repo.Config.OriginURL()
			must.EqOp(t, want, *have)
		}
	})

	t.Run("SetOffline", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		err := repo.Config.SetOffline(true)
		must.NoError(t, err)
		offline, err := repo.Config.IsOffline()
		must.NoError(t, err)
		must.True(t, offline)
		err = repo.Config.SetOffline(false)
		must.NoError(t, err)
		offline, err = repo.Config.IsOffline()
		must.NoError(t, err)
		must.False(t, offline)
	})
}
