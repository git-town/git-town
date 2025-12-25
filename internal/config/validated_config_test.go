package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/test/testruntime"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsMainOrPerennialBranch", func(t *testing.T) {
		t.Parallel()
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			},
		}
		tests := map[string]bool{
			"feature":     false,
			"main":        true,
			"perennial-1": true,
			"perennial-2": true,
			"perennial-3": false,
			"observed":    false,
			"parked":      false,
		}
		for give, want := range tests {
			have := config.IsMainOrPerennialBranch(gitdomain.NewLocalBranchName(give))
			fmt.Println(give)
			must.Eq(t, want, have)
		}
	})

	t.Run("Lineage", func(t *testing.T) {
		t.Parallel()
		repo := testruntime.CreateGitTown(t)
		repo.CreateFeatureBranch("feature1", "main")
		repo.CreateFeatureBranch("feature2", "main")
		repo.Config.Reload(repo.TestRunner)
		have := repo.Config.NormalConfig.Lineage
		want := configdomain.NewLineageWith(configdomain.LineageData{
			"feature1": "main",
			"feature2": "main",
		})
		must.Eq(t, want, have)
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			},
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})

	t.Run("RemoteURL", func(t *testing.T) {
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
			have, has := repo.Config.NormalConfig.RemoteURL(repo.TestRunner, gitdomain.RemoteOrigin).Get()
			must.True(t, has)
			must.EqOp(t, want, have)
		}
	})

	t.Run("RemovePerennialRoot", func(t *testing.T) {
		t.Parallel()
		contribution := gitdomain.NewLocalBranchName("contribution")
		feature1 := gitdomain.NewLocalBranchName("feature1")
		feature2 := gitdomain.NewLocalBranchName("feature2")
		main := gitdomain.NewLocalBranchName("main")
		observed := gitdomain.NewLocalBranchName("observed")
		perennial1 := gitdomain.NewLocalBranchName("perennial-1")
		perennial2 := gitdomain.NewLocalBranchName("perennial-2")
		perennialRegexOpt, err := configdomain.ParsePerennialRegex("peren*", "test")
		must.NoError(t, err)
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.LocalBranchNames{perennial1},
				PerennialRegex:    perennialRegexOpt,
			},
		}
		tests := map[*gitdomain.LocalBranchNames]gitdomain.LocalBranchNames{
			{main}:                           {},
			{perennial1, perennial2}:         {},
			{main, feature1, feature2}:       {feature1, feature2},
			{perennial1, feature1, feature2}: {feature1, feature2},
			{main, feature1, observed, contribution, feature2}: {feature1, observed, contribution, feature2},
			{main, perennial1, perennial2, feature1, feature2}: {feature1, feature2},
		}
		for give, want := range tests {
			have := config.RemovePerennials(*give)
			must.Eq(t, want, have)
		}
	})
}
