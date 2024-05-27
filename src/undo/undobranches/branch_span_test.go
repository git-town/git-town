package undobranches_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/undo/undobranches"
	"github.com/shoenig/test/must"
)

func TestBranchSpan(t *testing.T) {
	t.Parallel()

	t.Run("IsOmniChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.True(t, bs.IsOmniChange())
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.False(t, bs.IsOmniChange())
		})
	})

	t.Run("IsOmniRemove", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni remove", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			branchSpan := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
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
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
			}
			isOmniRemove, _, _ := bs.IsOmniRemove()
			must.False(t, isOmniRemove)
		})
	})

	t.Run("IsInconsistentChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an inconsistent change", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			}
			must.True(t, bs.IsInconsistentChange())
		})
		t.Run("no before-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no before-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no after-local", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no after-remote", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
	})

	t.Run("LocalAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("add a new local branch", func(t *testing.T) {
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isLocalAdded, afterBranchName, afterSHA := bs.LocalAdded()
			must.True(t, isLocalAdded)
			must.Eq(t, branch1, afterBranchName)
			must.Eq(t, sha1, afterSHA)
		})
		t.Run("add a local counterpart for an existing remote branch", func(t *testing.T) {
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
			}
			isLocalAdded, branchName, afterSHA := bs.LocalAdded()
			must.True(t, isLocalAdded)
			must.Eq(t, branch1, branchName)
			must.Eq(t, sha1, afterSHA)
		})
		t.Run("doesn't add anything", func(t *testing.T) {
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isLocalAdded, _, _ := bs.LocalAdded()
			must.False(t, isLocalAdded)
		})
	})

	t.Run("LocalChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changed a local branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.True(t, bs.LocalChanged())
		})
		t.Run("changed the local part of an omnibranch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.LocalChanged())
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.False(t, bs.LocalChanged())
		})
	})

	t.Run("LocalRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removed a local branch", func(t *testing.T) {
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isLocalRemoved, branchName, beforeSHA := bs.LocalRemoved()
			must.True(t, isLocalRemoved)
			must.Eq(t, branch1, branchName)
			must.Eq(t, sha1, beforeSHA)
		})
		t.Run("removed the local part of an omni branch", func(t *testing.T) {
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
			}
			isLocalRemoved, branchName, beforeSHA := bs.LocalRemoved()
			must.True(t, isLocalRemoved)
			must.Eq(t, branch1, branchName)
			must.Eq(t, sha1, beforeSHA)
		})
		t.Run("doesn't remove anything", func(t *testing.T) {
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isLocalRemoved, _, _ := bs.LocalRemoved()
			must.False(t, isLocalRemoved)
		})
	})

	t.Run("NoChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("no changes", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.NoChanges())
		})
		t.Run("has changes", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.False(t, bs.NoChanges())
		})
	})

	t.Run("RemoteAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("adds a remote-only branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				},
			}
			isRemoteAdded, addedRemoteBranchName, addedRemoteSHA := bs.RemoteAdded()
			must.True(t, isRemoteAdded)
			must.Eq(t, branch1, addedRemoteBranchName)
			must.Eq(t, sha1, addedRemoteSHA)
		})
		t.Run("adds the remote part for an existing local branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				},
			}
			isRemoteAdded, addedRemoteBranchName, addedRemoteSHA := bs.RemoteAdded()
			must.True(t, isRemoteAdded)
			must.Eq(t, branch1, addedRemoteBranchName)
			must.Eq(t, sha1, addedRemoteSHA)
		})
		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			isRemoteAdded, _, _ := bs.RemoteAdded()
			must.False(t, isRemoteAdded)
		})
	})

	t.Run("RemoteChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changes to a remote-only branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.True(t, bs.RemoteChanged())
		})
		t.Run("changes the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			must.True(t, bs.RemoteChanged())
		})
		t.Run("changes the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.RemoteChanged())
		})
	})

	t.Run("RemoteRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removing a remote-only branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isRemoteRemoved, remoteBranchName, beforeRemoteSHA := bs.RemoteRemoved()
			must.True(t, isRemoteRemoved)
			must.Eq(t, branch1, remoteBranchName)
			must.Eq(t, sha1, beforeRemoteSHA)
		})
		t.Run("removing the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				},
				After: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			isRemoteRemoved, remoteBranchName, beforeRemoteSHA := bs.RemoteRemoved()
			must.True(t, isRemoteRemoved)
			must.Eq(t, branch1, remoteBranchName)
			must.Eq(t, sha1, beforeRemoteSHA)
		})

		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				},
			}
			isRemoteRemoved, _, _ := bs.RemoteRemoved()
			must.False(t, isRemoteRemoved)
		})

		t.Run("upstream branch", func(t *testing.T) {
			t.Parallel()
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			isRemoteRemoved, _, _ := bs.RemoteRemoved()
			must.False(t, isRemoteRemoved)
		})
	})
}
