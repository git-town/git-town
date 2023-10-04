package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()

	t.Run("HasAllBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("has both branches", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.HasAllBranches())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasAllBranches())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.HasAllBranches())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasAllBranches())
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.True(t, give.HasLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.HasLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.HasLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasLocalBranch())
		})
	})

	t.Run("HasOnlyLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.True(t, give.HasOnlyLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasOnlyLocalBranch())
		})
	})

	t.Run("HasOnlyRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.HasOnlyRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasOnlyRemoteBranch())
		})
	})

	t.Run("HasRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.HasRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.HasRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.HasRemoteBranch())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.True(t, give.IsEmpty())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.IsEmpty())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.False(t, give.IsEmpty())
		})
	})

	t.Run("IsOmniBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			assert.True(t, give.IsOmniBranch())
		})
		t.Run("not an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("222222"),
			}
			assert.False(t, give.IsOmniBranch())
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			assert.False(t, give.IsOmniBranch())
		})
	})
}
