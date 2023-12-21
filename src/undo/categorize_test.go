package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/shoenig/test/must"
)

func TestCategorize(t *testing.T) {
	t.Parallel()

	t.Run("CategorizeInconsistentChanges", func(t *testing.T) {
		t.Parallel()
		give := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
				After: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
			},
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
				After: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1"),
		}
		havePerennials, haveFeatures := undo.CategorizeInconsistentChanges(give, branchTypes)
		wantPerennials := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
				After: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("perennial-1"),
					LocalSHA:   gitdomain.NewSHA("333333"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  gitdomain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("555555"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
				After: domain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("feature-1"),
					LocalSHA:   gitdomain.NewSHA("777777"),
					RemoteName: gitdomain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  gitdomain.NewSHA("888888"),
					SyncStatus: domain.SyncStatusNotInSync,
				},
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeLocalBranchChange", func(t *testing.T) {
		t.Parallel()
		give := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := undo.CategorizeLocalBranchChange(give, branchTypes)
		wantPerennials := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchChange", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := undo.CategorizeRemoteBranchChange(give, branchTypes)
		wantPerennials := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})

	t.Run("CategorizeRemoteBranchesSHAs", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("111111"),
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		havePerennials, haveFeatures := undo.CategorizeRemoteBranchesSHAs(give, branchTypes)
		wantPerennials := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"): gitdomain.NewSHA("111111"),
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
