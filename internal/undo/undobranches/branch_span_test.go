package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchSpan(t *testing.T) {
	t.Parallel()

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		t.Run("same branch name before and after", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
				}),
			}
			have := branchSpan.BranchNames()
			want := []gitdomain.BranchName{"branch", "origin/branch"}
			must.Eq(t, want, have)
		})

		t.Run("different branch name before and after", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-2"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-2")),
				}),
			}
			have := branchSpan.BranchNames()
			want := []gitdomain.BranchName{"branch-1", "branch-2", "origin/branch-1", "origin/branch-2"}
			must.Eq(t, want, have)
		})

		t.Run("all none", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: None[gitdomain.RemoteBranchName](),
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: None[gitdomain.RemoteBranchName](),
				}),
			}
			have := branchSpan.BranchNames()
			must.Len(t, 0, have)
		})
	})

	t.Run("IsInconsistentChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an inconsistent change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have, has := bs.InconsistentChange().Get()
			must.True(t, has)
			want := undodomain.InconsistentChange{
				After:  bs.After.GetOrPanic(),
				Before: bs.Before.GetOrPanic(),
			}
			must.Eq(t, want, have)
		})
		t.Run("no before-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			_, has := bs.InconsistentChange().Get()
			must.False(t, has)
		})
		t.Run("no before-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			_, has := bs.InconsistentChange().Get()
			must.False(t, has)
		})
		t.Run("no after-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			_, has := bs.InconsistentChange().Get()
			must.False(t, has)
		})
		t.Run("no after-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			_, has := bs.InconsistentChange().Get()
			must.False(t, has)
		})
	})

	t.Run("IsOmniChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni change", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha2}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha2),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			have, has := bs.OmniChange().Get()
			must.True(t, has)
			want := undobranches.LocalBranchChange{
				branch1: undodomain.Change[gitdomain.SHA]{
					Before: sha1,
					After:  sha2,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "222222"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			_, has := bs.OmniChange().Get()
			must.False(t, has)
		})
	})

	t.Run("IsOmniRemove", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni remove", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			have, has := branchSpan.OmniRemove().Get()
			must.True(t, has)
			want := undobranches.LocalBranchesSHAs{
				branch1: sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "333333"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha1}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			_, has := bs.OmniRemove().Get()
			must.False(t, has)
		})
	})

	t.Run("LocalAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("add a new local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			have, has := bs.LocalAdd().Get()
			must.True(t, has)
			must.EqOp(t, have, branch1)
		})
		t.Run("add a local counterpart for an existing remote branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			have, has := bs.LocalAdd().Get()
			must.True(t, has)
			must.EqOp(t, have, branch1)
		})
		t.Run("doesn't add anything", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After:  None[gitdomain.BranchInfo](),
			}
			_, has := bs.LocalAdd().Get()
			must.False(t, has)
		})
	})

	t.Run("LocalChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changed a local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha2}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			have, has := branchSpan.LocalChange().Get()
			must.True(t, has)
			want := undobranches.LocalBranchChange{
				branch1: undodomain.Change[gitdomain.SHA]{
					Before: sha1,
					After:  sha2,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changed the local part of an omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "222222"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have, has := branchSpan.LocalChange().Get()
			must.True(t, has)
			want := undobranches.LocalBranchChange{
				branch1: undodomain.Change[gitdomain.SHA]{
					Before: sha1,
					After:  sha2,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			_, has := branchSpan.LocalChange().Get()
			must.False(t, has)
		})
	})

	t.Run("LocalRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removed a local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			have, has := bs.LocalRemove().Get()
			must.True(t, has)
			want := undobranches.LocalBranchesSHAs{
				branch1: sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("removed the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: branch1, SHA: sha1}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have, has := bs.LocalRemove().Get()
			must.True(t, has)
			want := undobranches.LocalBranchesSHAs{
				branch1: sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("doesn't remove anything", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			_, has := bs.LocalRemove().Get()
			must.False(t, has)
		})
	})

	t.Run("RemoteAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("adds a remote-only branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have, has := bs.RemoteAdd().Get()
			must.True(t, has)
			must.Eq(t, branch1, have)
		})
		t.Run("adds the remote part for an existing local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha1}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha1}),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			have, has := bs.RemoteAdd().Get()
			must.True(t, has)
			must.Eq(t, branch1, have)
		})
		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			_, has := bs.RemoteAdd().Get()
			must.False(t, has)
		})
	})

	t.Run("RemoteChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changes to a remote-only branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha2),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have, has := branchSpan.RemoteChange().Get()
			must.True(t, has)
			want := undobranches.RemoteBranchChange{
				branch1: undodomain.Change[gitdomain.SHA]{
					Before: sha1,
					After:  sha2,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changes the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have, has := branchSpan.RemoteChange().Get()
			must.True(t, has)
			want := undobranches.RemoteBranchChange{
				branch1: undodomain.Change[gitdomain.SHA]{
					Before: sha1,
					After:  sha2,
				},
			}
			must.Eq(t, want, have)
		})
		t.Run("changes the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "111111"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: "222222"}),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			_, has := branchSpan.RemoteChange().Get()
			must.False(t, has)
		})
	})

	t.Run("RemoteRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removing a remote-only branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			have, has := bs.RemoteRemove().Get()
			must.True(t, has)
			want := undobranches.RemoteBranchesSHAs{
				branch1: sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("removing the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha1}),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      Some(gitdomain.BranchData{Name: "branch-1", SHA: sha1}),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			have, has := bs.RemoteRemove().Get()
			must.True(t, has)
			want := undobranches.RemoteBranchesSHAs{
				branch1: sha1,
			}
			must.Eq(t, want, have)
		})

		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			_, has := bs.RemoteRemove().Get()
			must.False(t, has)
		})

		t.Run("upstream branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					Local:      None[gitdomain.BranchData](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			_, has := bs.RemoteRemove().Get()
			must.False(t, has)
		})
	})
}
