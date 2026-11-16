package config_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v24/internal/config"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsMainOrPerennialBranch", func(t *testing.T) {
		t.Parallel()
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		tests := map[stringss.Trimmed]bool{
			"feature":     false,
			"main":        true,
			"observed":    false,
			"parked":      false,
			"perennial-1": true,
			"perennial-2": true,
			"perennial-3": false,
		}
		for give, want := range tests {
			have := config.IsMainOrPerennialBranch(gitdomain.LocalBranchNameOrPanic(give))
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
			NormalConfig: config.NormalConfig{
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})

	t.Run("RemovePerennialRoot", func(t *testing.T) {
		t.Parallel()
		contribution := gitdomain.LocalBranchName("contribution")
		feature1 := gitdomain.LocalBranchName("feature1")
		feature2 := gitdomain.LocalBranchName("feature2")
		main := gitdomain.LocalBranchName("main")
		observed := gitdomain.LocalBranchName("observed")
		perennial1 := gitdomain.LocalBranchName("perennial-1")
		perennial2 := gitdomain.LocalBranchName("perennial-2")
		perennialRegexOpt, err := configdomain.ParsePerennialRegex("peren*", "test")
		must.NoError(t, err)
		config := config.ValidatedConfig{
			NormalConfig: config.NormalConfig{
				PerennialBranches: gitdomain.LocalBranchNames{perennial1},
				PerennialRegex:    perennialRegexOpt,
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
			},
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
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
