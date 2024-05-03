package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
	"github.com/shoenig/test/must"
)

func TestCategorize(t *testing.T) {
	t.Parallel()

	t.Run("CategorizeInconsistentChanges", func(t *testing.T) {
		t.Parallel()
		give := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: configdomain.UnvalidatedConfig{ //nolint:exhaustruct
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1"),
			},
		}
		havePerennials, haveFeatures := undobranches.CategorizeInconsistentChanges(give, config)
		wantPerennials := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				},
				After: gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
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
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: configdomain.UnvalidatedConfig{ //nolint:exhaustruct
				PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
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
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: configdomain.UnvalidatedConfig{ //nolint:exhaustruct
				PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
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
		config := configdomain.ValidatedConfig{ //nolint:exhaustruct
			MainBranch: gitdomain.NewLocalBranchName("main"),
			UnvalidatedConfig: configdomain.UnvalidatedConfig{ //nolint:exhaustruct
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
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
