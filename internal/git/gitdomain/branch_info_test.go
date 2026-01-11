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
				Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			have, has := branchInfo.Local.Get()
			must.True(t, has)
			must.Eq(t, gitdomain.BranchData{Name: branch1, SHA: sha1}, have)
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(sha1),
			}
			have, has := branchInfo.Local.Get()
			must.True(t, has)
			must.Eq(t, gitdomain.BranchData{Name: branch1, SHA: sha1}, have)
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			_, has := branchInfo.Local.Get()
			must.False(t, has)
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			_, has := branchInfo.Local.Get()
			must.False(t, has)
		})
	})

	t.Run("GetLocalOrRemoteName", func(t *testing.T) {
		t.Parallel()
		t.Run("has local and remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch", SHA: "111111"}),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only local name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch", SHA: "111111"}),
				RemoteName: None[gitdomain.RemoteBranchName](),
			}
			have := branchInfo.GetLocalOrRemoteName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: "branch"}),
				RemoteName: gitdomain.NewRemoteBranchNameOption("origin/branch"),
			}
			have := branchInfo.GetLocalOrRemoteNameAsLocalName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only local name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch"}),
				RemoteName: None[gitdomain.RemoteBranchName](),
			}
			have := branchInfo.GetLocalOrRemoteNameAsLocalName()
			must.EqOp(t, "branch", have)
		})
		t.Run("has only remote name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
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
				Local:     Some(gitdomain.BranchData{SHA: "111111"}),
				RemoteSHA: Some(gitdomain.NewSHA("111111")),
			}
			have := branchInfo.GetLocalOrRemoteSHA()
			must.EqOp(t, "111111", have)
		})
		t.Run("has only local SHA", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:     Some(gitdomain.BranchData{SHA: "111111"}),
				RemoteSHA: None[gitdomain.SHA](),
			}
			have := branchInfo.GetLocalOrRemoteSHA()
			must.EqOp(t, "111111", have)
		})
		t.Run("has only remote SHA", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:     None[gitdomain.BranchData](),
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
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(branch1),
				RemoteSHA:  Some(sha1),
			}
			hasRemoteBranch, name, sha := branchInfo.GetRemote()
			must.True(t, hasRemoteBranch)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, sha)
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
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
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.True(t, give.HasOnlyLocalBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("has only a remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyLocalBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
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
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.True(t, give.HasOnlyRemoteBranch())
		})
		t.Run("has only a local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasOnlyRemoteBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.True(t, give.HasTrackingBranch())
		})
		t.Run("has local branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("has remote branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			must.False(t, give.HasTrackingBranch())
		})
		t.Run("is empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: branchName, SHA: "111111"}),
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
				Local:      Some(gitdomain.BranchData{Name: branchName, SHA: "111111"}),
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
				Local:      None[gitdomain.BranchData](),
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
				Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(sha1),
			}
			omni, has := give.OmniBranch().Get()
			must.True(t, has)
			must.EqOp(t, gitdomain.BranchData{Name: branch1, SHA: sha1}, omni)
		})
		t.Run("not an omnibranch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.SHA("222222")),
			}
			_, has := give.OmniBranch().Get()
			must.False(t, has)
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			_, has := give.OmniBranch().Get()
			must.False(t, has)
		})
	})
}
