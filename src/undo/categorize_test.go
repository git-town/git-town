package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
	"github.com/git-town/git-town/v11/src/undo/undobranches"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/shoenig/test/must"
)

func TestCategorize(t *testing.T) {
	t.Parallel()

	t.Run("CategorizeInconsistentChanges", func(t *testing.T) {
		t.Parallel()
		give := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
				After: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
			},
			undodomain.InconsistentChange{
				Before: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
				After: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
			},
		}
		branchTypes := syncdomain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1"),
		}
		havePerennials, haveFeatures := undobranches.CategorizeInconsistentChanges(give, branchTypes)
		wantPerennials := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
				After: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.InconsistentChanges{
			undodomain.InconsistentChange{
				Before: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
				After: undodomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
					SyncStatus: syncdomain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeLocalBranchChange", func(t *testing.T) {
		t.Parallel()
		give := undodomain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := syncdomain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := undobranches.CategorizeLocalBranchChange(give, branchTypes)
		wantPerennials := undodomain.LocalBranchChange{
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchChange", func(t *testing.T) {
		t.Parallel()
		give := undodomain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := syncdomain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchChange(give, branchTypes)
		wantPerennials := undodomain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchesSHAs", func(t *testing.T) {
		t.Parallel()
		give := undodomain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("111111"),
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		branchTypes := syncdomain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		havePerennials, haveFeatures := undobranches.CategorizeRemoteBranchesSHAs(give, branchTypes)
		wantPerennials := undodomain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := undodomain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"): gitdomain.NewSHA("111111"),
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
