package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is a local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(branch1),
				LocalSHA:   Some(sha1),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			hasLocalBranch, name, sha := branchInfo.HasLocalBranch()
			must.True(t, hasLocalBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(branch1),
				LocalSHA:   Some(sha1),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(sha1),
			}
			hasLocalBranch, name, sha := branchInfo.HasLocalBranch()
			must.True(t, hasLocalBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			hasLocalBranch, _, _ := branchInfo.HasLocalBranch()
			must.False(t, hasLocalBranch)
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			hasLocalBranch, _, _ := branchInfo.HasLocalBranch()
			must.False(t, hasLocalBranch)
		})
	})

	t.Run("HasOnlyLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.True(t, give.HasOnlyLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
	})

	t.Run("HasOnlyRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.True(t, give.HasOnlyRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
	})

	t.Run("HasRemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(branch1),
				RemoteSHA:  Some(sha1),
			}
			hasRemoteBranch, name, sha := branchInfo.HasRemoteBranch()
			must.True(t, hasRemoteBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			hasRemoteBranch, name, sha := branchInfo.HasRemoteBranch()
			must.True(t, hasRemoteBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			hasRemoteBranch, _, _ := branchInfo.HasRemoteBranch()
			must.False(t, hasRemoteBranch)
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			hasRemoteBranch, _, _ := branchInfo.HasRemoteBranch()
			must.False(t, hasRemoteBranch)
		})
	})

	t.Run("HasTrackingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has both branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.True(t, give.HasTrackingBranch())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasTrackingBranch())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.True(t, give.IsEmpty())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.IsEmpty())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.IsEmpty())
		})
	})

	t.Run("IsOmniBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			give := gitdomain.BranchInfo{
				LocalName:  Some(branch1),
				LocalSHA:   Some(sha1),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(sha1),
			}
			isOmni, name, sha := give.IsOmniBranch()
			must.True(t, isOmni)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("not an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			}
			isOmni, _, _ := give.IsOmniBranch()
			must.False(t, isOmni)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			isOmni, _, _ := give.IsOmniBranch()
			must.False(t, isOmni)
		})
	})
}
