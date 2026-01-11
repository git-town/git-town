package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()

	t.Run("BranchIsActiveInAnotherWorktree", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is active in another worktree", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-2")
			must.True(t, have)
		})
		t.Run("branch is local but not active in another worktree", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-1")
			must.False(t, have)
		})
		t.Run("branch is remote", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-1")
			must.False(t, have)
		})
		t.Run("branch doesn't exist", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("zonk")
			must.False(t, have)
		})
	})

	t.Run("BranchesDeletedAtRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("empty BranchInfos", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{}
			have := branchInfos.BranchesDeletedAtRemote()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no branches deleted at remote", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
				},
			}
			have := branchInfos.BranchesDeletedAtRemote()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("multiple branches deleted at remote", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-3"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-4"}),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
				{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-5")),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
			}
			have := branchInfos.BranchesDeletedAtRemote()
			want := gitdomain.NewLocalBranchNames("branch-2", "branch-4", "branch-5")
			must.Eq(t, want, have)
		})
		t.Run("mixed statuses including deleted at remote", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "local-only"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "up-to-date"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "deleted-at-remote-1"}),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "ahead"}),
					SyncStatus: gitdomain.SyncStatusAhead,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "other-worktree"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "deleted-at-remote-2"}),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
			}
			have := branchInfos.BranchesDeletedAtRemote()
			want := gitdomain.NewLocalBranchNames("deleted-at-remote-1", "deleted-at-remote-2")
			must.Eq(t, want, have)
		})
	})

	t.Run("BranchesInOtherWorktrees", func(t *testing.T) {
		t.Parallel()
		t.Run("empty BranchInfos", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("no branches in other worktrees", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
				},
			}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("one branch in another worktree with local name", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
			}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.NewLocalBranchNames("branch-2")
			must.Eq(t, want, have)
		})
		t.Run("one branch in another worktree with remote name only", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
			}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.NewLocalBranchNames("branch-2")
			must.Eq(t, want, have)
		})
		t.Run("multiple branches in other worktrees", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-3"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-3")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "branch-4"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
				{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-5")),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
			}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.NewLocalBranchNames("branch-2", "branch-4", "branch-5")
			must.Eq(t, want, have)
		})
		t.Run("mixed statuses including other worktrees", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					Local:      Some(gitdomain.BranchData{Name: "local-only"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "up-to-date"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "other-worktree-1"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "ahead"}),
					SyncStatus: gitdomain.SyncStatusAhead,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "other-worktree-2"}),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
				{
					Local:      Some(gitdomain.BranchData{Name: "deleted-at-remote"}),
					SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				},
			}
			have := branchInfos.BranchesInOtherWorktrees()
			want := gitdomain.NewLocalBranchNames("other-worktree-1", "other-worktree-2")
			must.Eq(t, want, have)
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/two")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bs := gitdomain.BranchInfos{branch}
			have, has := bs.FindByRemoteName("origin/two").Get()
			must.True(t, has)
			must.Eq(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "kg/one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}}
			have := bs.FindByRemoteName("kg/one")
			must.True(t, have.IsNone())
		})
	})

	t.Run("FindLocalOrRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has local name", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch1info := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: branch1, SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(branch1)
			must.Eq(t, MutableSome(&branch1info), have)
		})
		t.Run("has remote name", func(t *testing.T) {
			t.Parallel()
			branch1info := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(gitdomain.NewLocalBranchName("branch-1"))
			must.Eq(t, MutableSome(&branch1info), have)
		})
		t.Run("no match", func(t *testing.T) {
			t.Parallel()
			branch1info := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(gitdomain.NewLocalBranchName("zonk"))
			must.Eq(t, MutableNone[gitdomain.BranchInfo](), have)
		})
	})

	t.Run("FindMatchingRecord", func(t *testing.T) {
		t.Parallel()
		t.Run("has matching local name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			give := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			have := bis.FindMatchingRecord(give)
			want := MutableSome(&bis[0])
			must.Eq(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			give := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			have := bis.FindMatchingRecord(give)
			want := MutableSome(&bis[0])
			must.Eq(t, want, have)
		})
	})

	t.Run("FindRemoteNameMatchingLocal", func(t *testing.T) {
		t.Parallel()
		t.Run("has a remote branch matching the local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "other", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/target")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bs := gitdomain.BranchInfos{branch}
			have, has := bs.FindRemoteNameMatchingLocal("target").Get()
			must.True(t, has)
			must.Eq(t, &branch, have)
		})
		t.Run("has a remote-only branch matching the local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/target")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bs := gitdomain.BranchInfos{branch}
			have, has := bs.FindRemoteNameMatchingLocal("target").Get()
			must.True(t, has)
			must.Eq(t, &branch, have)
		})
		t.Run("has a remote branch with different local branch name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/two")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			have := bs.FindRemoteNameMatchingLocal("target")
			must.True(t, have.IsNone())
		})
		t.Run("has a branch without remote name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "target", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			have := bs.FindRemoteNameMatchingLocal(gitdomain.NewLocalBranchName("target"))
			must.True(t, have.IsNone())
		})
		t.Run("empty BranchInfos", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{}
			have := bs.FindRemoteNameMatchingLocal(gitdomain.NewLocalBranchName("target"))
			must.True(t, have.IsNone())
		})
		t.Run("multiple branches, one matching", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/other")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			branch2 := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "222222"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/target")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			}
			branch3 := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "three", SHA: "333333"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/another")),
				RemoteSHA:  Some(gitdomain.NewSHA("333333")),
			}
			bs := gitdomain.BranchInfos{branch1, branch2, branch3}
			have, has := bs.FindRemoteNameMatchingLocal("target").Get()
			must.True(t, has)
			must.Eq(t, &branch2, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.True(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with a matching remote", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.False(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "up-to-date", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "ahead", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "behind", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "local-only", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "deleted-at-remote", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.LocalBranches().NamesLocalBranches()
		want := gitdomain.NewLocalBranchNames("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "up-to-date", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "ahead", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "behind", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "local-only", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "deleted-at-remote", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().NamesLocalBranches()
		want := gitdomain.NewLocalBranchNames("deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("FindByLocalName", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch with matching name", func(t *testing.T) {
			t.Parallel()
			branchInfo := gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			branchInfos := gitdomain.BranchInfos{branchInfo}
			have, has := branchInfos.FindByLocalName("one").Get()
			must.True(t, has)
			must.Eq(t, &branchInfo, have)
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("kg/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			have := branchInfos.FindByLocalName(gitdomain.NewLocalBranchName("kg/one"))
			must.True(t, have.IsNone())
		})
	})

	t.Run("NamesAllBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/three")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
		}
		have := bs.NamesAllBranches()
		want := gitdomain.NewLocalBranchNames("one", "two", "three")
		must.Eq(t, want, have)
	})

	t.Run("NamesLocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      None[gitdomain.BranchData](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/three")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
		}
		have := bs.NamesLocalBranches()
		want := gitdomain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			have := bs.Remove(gitdomain.NewLocalBranchName("two"))
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "three", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "four", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have, nonExisting := bs.Select("one", "three")
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "three", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.Eq(t, want, have)
		must.Eq(t, nonExisting, gitdomain.LocalBranchNames(nil))
	})

	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.Remove(gitdomain.NewLocalBranchName("zonk"))
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "one", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				Local:      Some(gitdomain.BranchData{Name: "two", SHA: "111111"}),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.Eq(t, want, have)
	})
}
