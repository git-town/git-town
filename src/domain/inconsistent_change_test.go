package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestInconsistentChange(t *testing.T) {
	t.Parallel()
	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		ics := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-1"),
					LocalSHA:   domain.NewSHA("111111"),
					RemoteName: domain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusBehind,
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-1"),
					LocalSHA:   domain.NewSHA("333333"),
					RemoteName: domain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusBehind,
				},
			},
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-1"),
					LocalSHA:   domain.NewSHA("555555"),
					RemoteName: domain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  domain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusBehind,
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-1"),
					LocalSHA:   domain.NewSHA("777777"),
					RemoteName: domain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  domain.NewSHA("888888"),
					SyncStatus: domain.SyncStatusBehind,
				},
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-1"),
		}
		havePerennials, haveFeatures := ics.Categorize(branchTypes)
		wantPerennials := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-1"),
					LocalSHA:   domain.NewSHA("111111"),
					RemoteName: domain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusBehind,
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("perennial-1"),
					LocalSHA:   domain.NewSHA("333333"),
					RemoteName: domain.NewRemoteBranchName("origin/perennial-1"),
					RemoteSHA:  domain.NewSHA("444444"),
					SyncStatus: domain.SyncStatusBehind,
				},
			},
		}
		assert.Equal(t, wantPerennials, havePerennials)
		wantFeatures := domain.InconsistentChanges{
			domain.InconsistentChange{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-1"),
					LocalSHA:   domain.NewSHA("555555"),
					RemoteName: domain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  domain.NewSHA("666666"),
					SyncStatus: domain.SyncStatusBehind,
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("feature-1"),
					LocalSHA:   domain.NewSHA("777777"),
					RemoteName: domain.NewRemoteBranchName("origin/feature-1"),
					RemoteSHA:  domain.NewSHA("888888"),
					SyncStatus: domain.SyncStatusBehind,
				},
			},
		}
		assert.Equal(t, wantFeatures, haveFeatures)
	})
}
