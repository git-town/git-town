package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.True(t, give.HasLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.HasLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.HasLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasLocalBranch())
		})
	})

	t.Run("HasOnlyLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.True(t, give.HasOnlyLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
	})

	t.Run("HasOnlyRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.HasOnlyRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
	})

	t.Run("HasRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.HasRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.HasRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasRemoteBranch())
		})
	})

	t.Run("HasTrackingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has both branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.HasTrackingBranch())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.HasTrackingBranch())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.True(t, give.IsEmpty())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.IsEmpty())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.False(t, give.IsEmpty())
		})
	})

	t.Run("IsOmniBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			}
			must.True(t, give.IsOmniBranch())
		})
		t.Run("not an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("222222"),
			}
			must.False(t, give.IsOmniBranch())
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			must.False(t, give.IsOmniBranch())
		})
	})
}
