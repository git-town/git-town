package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestInconsistentChange(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
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
		havePerennials, haveFeatures := give.Categorize(branchTypes)
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
}
