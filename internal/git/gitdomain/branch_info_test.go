package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()

	t.Run("GetLocal", func(t *testing.T) {
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
			hasLocalBranch, name, sha := branchInfo.GetLocal()
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
			hasLocalBranch, name, sha := branchInfo.GetLocal()
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
			hasLocalBranch, _, _ := branchInfo.GetLocal()
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
			hasLocalBranch, _, _ := branchInfo.GetLocal()
			must.False(t, hasLocalBranch)
		})
	})

	t.Run("GetLocalOrRemoteName", func(t *testing.T) {
		t.Parallel()
		t.Run("has local and remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only local name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
				RemoteName: None[gitdomain.RemoteBranchName](),
			}
			have := branchInfo.GetLocalOrRemoteName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteName()
			must.EqOp(t, "origin/branch", have)
		})
	})

	t.Run("GetLocalOrRemoteNameAsLocalName", func(t *testing.T) {
		t.Parallel()
		t.Run("has local and remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteNameAsLocalName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only local name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
				RemoteName: None[gitdomain.RemoteBranchName](),
			}
			have := branchInfo.GetLocalOrRemoteNameAsLocalName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteNameAsLocalName()
			must.EqOp(t, "branch", have)
		})
	})

	t.Run("GetLocalOrRemoteSHA", func(t *testing.T) {
		t.Parallel()
		t.Run("has local and remote SHA", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalSHA:  Some(gitdomain.NewSHA("111111")),
				RemoteSHA: Some(gitdomain.NewSHA("111111")),
			}
			have := branchInfo.GetLocalOrRemoteSHA()
			must.EqOp(t, "111111", have)
		})
		t.Run("has only local SHA", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalSHA:  Some(gitdomain.NewSHA("111111")),
				RemoteSHA: None[gitdomain.SHA](),
			}
			have := branchInfo.GetLocalOrRemoteSHA()
			must.EqOp(t, "111111", have)
		})
		t.Run("has only remote SHA", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalSHA:  None[gitdomain.SHA](),
				RemoteSHA: Some(gitdomain.NewSHA("111111")),
			}
			have := branchInfo.GetLocalOrRemoteSHA()
			must.EqOp(t, "111111", have)
		})
	})

	t.Run("GetRemoteBranch", func(t *testing.T) {
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
			hasRemoteBranch, name, sha := branchInfo.GetRemote()
			must.True(t, hasRemoteBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			hasRemoteBranch, name, sha := branchInfo.GetRemote()
			must.True(t, hasRemoteBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			hasRemoteBranch, _, _ := branchInfo.GetRemote()
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
			hasRemoteBranch, _, _ := branchInfo.GetRemote()
			must.False(t, hasRemoteBranch)
		})
	})

	t.Run("HasOnlyLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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

	t.Run("HasTrackingBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has both branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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

	t.Run("IsLocalOnlyBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("is indeed a local branch", func(t *testing.T) {
			t.Parallel()
			branchName := gitdomain.NewLocalBranchName("foo")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(branchName),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
			}
			isLocal, haveBranchName := branchInfo.IsLocalOnlyBranch()
			must.True(t, isLocal)
			must.Eq(t, branchName, haveBranchName)
		})
		t.Run("has a tracking branch", func(t *testing.T) {
			t.Parallel()
			branchName := gitdomain.NewLocalBranchName("foo")
			branchInfo := gitdomain.BranchInfo{
				LocalName:  Some(branchName),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/foo")),
				RemoteSHA:  Some(gitdomain.SHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
			}
			isLocal, haveBranchName := branchInfo.IsLocalOnlyBranch()
			must.False(t, isLocal)
			must.Eq(t, branchName, haveBranchName)
		})
		t.Run("remote-only branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/foo")),
				RemoteSHA:  Some(gitdomain.SHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
			}
			isLocal, haveBranchName := branchInfo.IsLocalOnlyBranch()
			must.False(t, isLocal)
			must.Eq(t, "", haveBranchName)
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
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
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
