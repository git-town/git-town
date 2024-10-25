package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/undo/undobranches"
	"github.com/git-town/git-town/v16/internal/undo/undodomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestCategorize(t *testing.T) {
	t.Parallel()

	t.Run("CategorizeInconsistentChanges", func(t *testing.T) {
		t.Parallel()
		give := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-1")),
					LocalSHA:   Some(gitdomain.NewSHA("555555")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-1")),
					LocalSHA:   Some(gitdomain.NewSHA("777777")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("888888")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: gitdomain.NewLocalBranchName("main"),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1"),
				},
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeInconsistentChanges(give, config)
		wantPerennials := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/perennial-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("perennial-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
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
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-1")),
					LocalSHA:   Some(gitdomain.NewSHA("555555")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/feature-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("666666")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("feature-1")),
					LocalSHA:   Some(gitdomain.NewSHA("777777")),
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
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: gitdomain.NewLocalBranchName("main"),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
				},
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeLocalBranchChange(give, config)
		wantPerennials := undobranches.LocalBranchChange{
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchChange", func(t *testing.T) {
		t.Parallel()
		give := undobranches.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: gitdomain.NewLocalBranchName("main"),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
				},
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchChange(give, config)
		wantPerennials := undobranches.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchesSHAs", func(t *testing.T) {
		t.Parallel()
		give := undobranches.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("111111"),
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		config := config.ValidatedConfig{
			ValidatedConfigData: configdomain.ValidatedConfigData{
				MainBranch: gitdomain.NewLocalBranchName("main"),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
				},
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchesSHAs(give, config)
		wantPerennials := undobranches.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undobranches.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"): gitdomain.NewSHA("111111"),
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
