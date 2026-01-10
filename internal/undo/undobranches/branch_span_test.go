package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undobranches"
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
					LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch")),
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch"),
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
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-2"),
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
					LocalName:  None[gitdomain.LocalBranchName](),
					RemoteName: None[gitdomain.RemoteBranchName](),
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
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
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			isInconsistentChange, before, after := bs.IsInconsistentChange()
			must.True(t, isInconsistentChange)
			must.Eq(t, bs.Before.GetOrPanic(), before)
			must.Eq(t, bs.After.GetOrPanic(), after)
		})
		t.Run("no before-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			isInconsistentChange, _, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
		})
		t.Run("no before-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			isInconsistentChange, _, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
		})
		t.Run("no after-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			isInconsistentChange, _, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
		})
		t.Run("no after-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			isInconsistentChange, _, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
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
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha2),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha2),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			isOmni, name, beforeSHA, afterSHA := bs.IsOmniChange()
			must.True(t, isOmni)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, beforeSHA)
			must.EqOp(t, sha2, afterSHA)
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			isOmni, _, _, _ := bs.IsOmniChange()
			must.False(t, isOmni)
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
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			isOmniRemove, beforeBranchName, beforeSHA := branchSpan.IsOmniRemove()
			must.True(t, isOmniRemove)
			must.Eq(t, branch1, beforeBranchName)
			must.Eq(t, sha1, beforeSHA)
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			sha1 := gitdomain.NewSHA("111111")
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			isOmniRemove, _, _ := bs.IsOmniRemove()
			must.False(t, isOmniRemove)
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
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			isLocalAdded, afterBranchName, afterSHA := bs.LocalAdded()
			must.True(t, isLocalAdded)
			must.Eq(t, branch1, afterBranchName)
			must.Eq(t, sha1, afterSHA)
		})
		t.Run("add a local counterpart for an existing remote branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			isLocalAdded, branchName, afterSHA := bs.LocalAdded()
			must.True(t, isLocalAdded)
			must.Eq(t, branch1, branchName)
			must.Eq(t, sha1, afterSHA)
		})
		t.Run("doesn't add anything", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: None[gitdomain.BranchInfo](),
				After:  None[gitdomain.BranchInfo](),
			}
			isLocalAdded, _, _ := bs.LocalAdded()
			must.False(t, isLocalAdded)
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
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(sha2),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			have := branchSpan.LocalChanged()
			want := undobranches.LocalChangedResult{
				IsChanged: true,
				Name:      branch1,
				SHABefore: sha1,
				SHAAfter:  sha2,
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
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have := branchSpan.LocalChanged()
			want := undobranches.LocalChangedResult{
				IsChanged: true,
				Name:      branch1,
				SHABefore: sha1,
				SHAAfter:  sha2,
			}
			must.Eq(t, want, have)
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have := branchSpan.LocalChanged()
			want := undobranches.LocalChangedResult{
				IsChanged: false,
				Name:      "branch-1",
				SHABefore: "111111",
				SHAAfter:  "111111",
			}
			must.Eq(t, want, have)
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
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			have := bs.LocalRemoved()
			want := undobranches.LocalRemovedResult{
				IsRemoved: true,
				Name:      branch1,
				SHA:       sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("removed the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have := bs.LocalRemoved()
			want := undobranches.LocalRemovedResult{
				IsRemoved: true,
				Name:      branch1,
				SHA:       sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("doesn't remove anything", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			have := bs.LocalRemoved()
			want := undobranches.LocalRemovedResult{
				IsRemoved: false,
				Name:      "branch-1",
				SHA:       "111111",
			}
			must.Eq(t, want, have)
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
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have := bs.RemoteAdded()
			want := undobranches.RemoteAddedResult{
				IsAdded: true,
				Name:    branch1,
				SHA:     sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("adds the remote part for an existing local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(sha1),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(sha1),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
			}
			have := bs.RemoteAdded()
			want := undobranches.RemoteAddedResult{
				IsAdded: true,
				Name:    branch1,
				SHA:     sha1,
			}
			must.Eq(t, want, have)
		})
		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have := bs.RemoteAdded()
			want := undobranches.RemoteAddedResult{
				IsAdded: false,
				Name:    "origin/branch-1",
				SHA:     "222222",
			}
			must.Eq(t, want, have)
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
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha2),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			have := branchSpan.RemoteChanged()
			want := undobranches.RemoteChangedResult{
				IsChanged: true,
				Name:      branch1,
				SHABefore: sha1,
				SHAAfter:  sha2,
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
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have := branchSpan.RemoteChanged()
			want := undobranches.RemoteChangedResult{
				IsChanged: true,
				Name:      branch1,
				SHABefore: sha1,
				SHAAfter:  sha2,
			}
			must.Eq(t, want, have)
		})
		t.Run("changes the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
				}),
			}
			have := branchSpan.RemoteChanged()
			want := undobranches.RemoteChangedResult{
				IsChanged: false,
				Name:      "origin/branch-1",
				SHABefore: "111111",
				SHAAfter:  "111111",
			}
			must.Eq(t, want, have)
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
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: None[gitdomain.BranchInfo](),
			}
			remoteRemoved := bs.RemoteRemoved()
			must.True(t, remoteRemoved.IsRemoved)
			must.Eq(t, branch1, remoteRemoved.Name)
			must.Eq(t, sha1, remoteRemoved.SHA)
		})
		t.Run("removing the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(sha1),
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(sha1),
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				}),
			}
			remoteRemoved := bs.RemoteRemoved()
			must.True(t, remoteRemoved.IsRemoved)
			must.Eq(t, branch1, remoteRemoved.Name)
			must.Eq(t, sha1, remoteRemoved.SHA)
		})

		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			remoteRemoved := bs.RemoteRemoved()
			must.False(t, remoteRemoved.IsRemoved)
		})

		t.Run("upstream branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
				}),
			}
			remoteRemoved := bs.RemoteRemoved()
			must.False(t, remoteRemoved.IsRemoved)
		})
	})
}
