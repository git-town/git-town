package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()
	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			bi := domain.BranchInfo{
				LocalName:  domain.LocalBranchName{},
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			}
			assert.True(t, bi.IsEmpty())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			bi := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			}
			assert.False(t, bi.IsEmpty())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			bi := domain.BranchInfo{
				LocalName:  domain.LocalBranchName{},
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, bi.IsEmpty())
		})
	})

	t.Run("IsOmniBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omnibranch", func(t *testing.T) {
			t.Parallel()
			bi := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, bi.IsOmniBranch())
		})
		t.Run("not an omnibranch", func(t *testing.T) {
			t.Parallel()
			bi := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("222222"),
			}
			assert.False(t, bi.IsOmniBranch())
		})
	})
}
