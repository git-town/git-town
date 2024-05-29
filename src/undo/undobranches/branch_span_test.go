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
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			bs := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha2),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha2),
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
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
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
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				},
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
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				}),
			}
			isInconsistentChange, after := bs.IsInconsistentChange()
			must.True(t, isInconsistentChange)
			must.Eq(t, bs.After.GetOrPanic(), after)
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				}),
			}
			isInconsistentChange, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				}),
			}
			isInconsistentChange, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("444444")),
				}),
			}
			isInconsistentChange, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("333333")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				}),
			}
			isInconsistentChange, _ := bs.IsInconsistentChange()
			must.False(t, isInconsistentChange)
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				}),
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
				After: None[gitdomain.BranchInfo](),
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
				Before: gitdomain.BranchInfo{
					LocalName:  Some(branch1),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha2),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				}),
			}
			localChanged, name, beforeSHA, afterSHA := branchSpan.LocalChanged()
			must.True(t, localChanged)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, beforeSHA)
			must.EqOp(t, sha2, afterSHA)
		})
		t.Run("changed the local part of an omnibranch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				}),
			}
			localChanged, name, beforeSHA, afterSHA := branchSpan.LocalChanged()
			must.True(t, localChanged)
			must.EqOp(t, branch1, name)
			must.EqOp(t, sha1, beforeSHA)
			must.EqOp(t, sha2, afterSHA)
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				}),
			}
			localChanged, _, _, _ := branchSpan.LocalChanged()
			must.False(t, localChanged)
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
				After: None[gitdomain.BranchInfo](),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(sha1),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				}),
			}
			isRemoteAdded, _, _ := bs.RemoteAdded()
			must.False(t, isRemoteAdded)
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
				Before: gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha1),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(branch1),
					RemoteSHA:  Some(sha2),
				}),
			}
			remoteChanged, branch, beforeSHA, afterSHA := branchSpan.RemoteChanged()
			must.True(t, remoteChanged)
			must.Eq(t, branch1, branch)
			must.Eq(t, sha1, beforeSHA)
			must.Eq(t, sha2, afterSHA)
		})
		t.Run("changes the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewRemoteBranchName("origin/branch-1")
			sha1 := gitdomain.NewSHA("111111")
			sha2 := gitdomain.NewSHA("222222")
			branchSpan := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				}),
			}
			remoteChanged, branch, beforeSHA, afterSHA := branchSpan.RemoteChanged()
			must.True(t, remoteChanged)
			must.Eq(t, branch1, branch)
			must.Eq(t, sha1, beforeSHA)
			must.Eq(t, sha2, afterSHA)
		})
		t.Run("changes the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			branchSpan := undobranches.BranchSpan{
				Before: gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("222222")),
					SyncStatus: gitdomain.SyncStatusNotInSync,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				}),
			}
			remoteChanged, _, _, _ := branchSpan.RemoteChanged()
			must.False(t, remoteChanged)
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
				After: None[gitdomain.BranchInfo](),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(sha1),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("222222")),
				}),
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
				After: Some(gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("upstream/main")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				}),
			}
			isRemoteRemoved, _, _ := bs.RemoteRemoved()
			must.False(t, isRemoteRemoved)
		})
	})
}
