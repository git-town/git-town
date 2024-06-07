package config_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("Author", func(t *testing.T) {
		t.Parallel()
		conf := config.ValidatedConfig{
			Config: configdomain.ValidatedConfig{
				GitUserName:  configdomain.GitUserName("name"),
				GitUserEmail: configdomain.GitUserEmail("email"),
			},
		}
		have := conf.Author()
		want := gitdomain.Author("name <email>")
		must.EqOp(t, want, have)
	})

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		repo.CreateFeatureBranch(gitdomain.NewLocalBranchName("feature1"))
		repo.CreateFeatureBranch(gitdomain.NewLocalBranchName("feature2"))
		repo.Config.Reload()
		have := repo.Config.Config.Lineage
		want := configdomain.NewLineage()
		want.AddParent(gitdomain.NewLocalBranchName("feature1"), gitdomain.NewLocalBranchName("main"))
		want.AddParent(gitdomain.NewLocalBranchName("feature2"), gitdomain.NewLocalBranchName("main"))
		must.Eq(t, want, have)
	})

	t.Run("OriginURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]giturl.Parts{
			"http://github.com/organization/repository":                     {Host: "github.com", Org: "organization", Repo: "repository", User: None[string]()},
			"http://github.com/organization/repository.git":                 {Host: "github.com", Org: "organization", Repo: "repository", User: None[string]()},
			"https://github.com/organization/repository":                    {Host: "github.com", Org: "organization", Repo: "repository", User: None[string]()},
			"https://github.com/organization/repository.git":                {Host: "github.com", Org: "organization", Repo: "repository", User: None[string]()},
			"https://sub.domain.customhost.com/organization/repository":     {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository", User: None[string]()},
			"https://sub.domain.customhost.com/organization/repository.git": {Host: "sub.domain.customhost.com", Org: "organization", Repo: "repository", User: None[string]()},
		}
		for give, want := range tests {
			repo := testruntime.CreateGitTown(t)
			os.Setenv("GIT_TOWN_REMOTE", give)
			defer os.Unsetenv("GIT_TOWN_REMOTE")
			have, has := repo.Config.OriginURL().Get()
			must.True(t, has)
			must.EqOp(t, want, have)
		}
	})

	t.Run("Reload", func(t *testing.T) {
		t.Parallel()
		t.Run("lineage changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			branch := gitdomain.NewLocalBranchName("branch-1")
			repo.CreateFeatureBranch(branch)
			repo.Config.Reload()
			want := configdomain.NewLineage()
			want.AddParent(branch, gitdomain.NewLocalBranchName("main"))
			must.Eq(t, want, repo.Config.Config.Lineage)
		})
	})
}
