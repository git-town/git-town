package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestCategorize(t *testing.T) {
	t.Parallel()

	t.Run("CategorizeInconsistentChanges", func(t *testing.T) {
		t.Parallel()
		give := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-1", SHA: "555555"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-1", SHA: "777777"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("888888")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1"),
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeInconsistentChanges(give, config)
		wantPerennials := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "perennial-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-1", SHA: "555555"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "feature-1", SHA: "777777"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("888888")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeLocalBranchChange", func(t *testing.T) {
		t.Parallel()
		give := undobranches.LocalBranchChange{
			"branch-1": {
				Before: "111111",
				After:  "222222",
			},
			"dev": {
				Before: "333333",
				After:  "444444",
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeLocalBranchChange(give, config)
		wantPerennials := undobranches.LocalBranchChange{
			"dev": {
				Before: "333333",
				After:  "444444",
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.LocalBranchChange{
			"branch-1": {
				Before: "111111",
				After:  "222222",
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchChange", func(t *testing.T) {
		t.Parallel()
		give := undobranches.RemoteBranchChange{
			"origin/branch-1": {
				Before: "111111",
				After:  "222222",
			},
			"origin/dev": {
				Before: "333333",
				After:  "444444",
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchChange(give, config)
		wantPerennials := undobranches.RemoteBranchChange{
			"origin/dev": {
				Before: "333333",
				After:  "444444",
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.RemoteBranchChange{
			"origin/branch-1": {
				Before: "111111",
				After:  "222222",
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchesSHAs", func(t *testing.T) {
		t.Parallel()
		give := undobranches.RemoteBranchesSHAs{
			"origin/feature-branch":   "111111",
			"origin/perennial-branch": "222222",
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: "main",
			},
			NormalConfig: config.NormalConfig{
				UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchesSHAs(give, config)
		wantPerennials := undobranches.RemoteBranchesSHAs{
			"origin/perennial-branch": "222222",
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.RemoteBranchesSHAs{
			"origin/feature-branch": "111111",
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
