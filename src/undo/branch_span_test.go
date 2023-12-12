package undo_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/undo"
	"github.com/shoenig/test/must"
)

func TestBranchSpan(t *testing.T) {
	t.Parallel()

	t.Run("IsOmniChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.True(t, bs.IsOmniChange())
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.False(t, bs.IsOmniChange())
		})
	})

	t.Run("IsOmniRemove", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni remove", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.IsOmniRemove())
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.IsOmniRemove())
		})
	})

	t.Run("IsInconsistentChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an inconsistent change", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			}
			must.True(t, bs.IsInconsistentChange())
		})
		t.Run("no before-local", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no before-remote", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no after-local", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("444444"),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
		t.Run("no after-remote", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusNotInSync,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.IsInconsistentChange())
		})
	})

	t.Run("LocalAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("add a new local branch", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.LocalAdded())
		})
		t.Run("add a local counterpart for an existing remote branch", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.LocalAdded())
		})
		t.Run("doesn't add anything", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.LocalAdded())
		})
	})

	t.Run("LocalChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changed a local branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.LocalChanged())
		})
		t.Run("changed the local part of an omnibranch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.LocalChanged())
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.False(t, bs.LocalChanged())
		})
	})

	t.Run("LocalRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removed a local branch", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.LocalRemoved())
		})
		t.Run("removed the local part of an omni branch", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.LocalRemoved())
		})
		t.Run("doesn't remove anything", func(t *testing.T) {
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.LocalAdded())
		})
	})

	t.Run("NoChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("no changes", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.NoChanges())
		})
		t.Run("has changes", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.False(t, bs.NoChanges())
		})
	})

	t.Run("RemoteAdded", func(t *testing.T) {
		t.Parallel()
		t.Run("adds a remote-only branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.RemoteAdded())
		})
		t.Run("adds the remote part for an existing local branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.True(t, bs.RemoteAdded())
		})
		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.False(t, bs.RemoteAdded())
		})
	})

	t.Run("RemoteChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("changes a remote-only branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.True(t, bs.RemoteChanged())
		})
		t.Run("changes the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.True(t, bs.RemoteChanged())
		})
		t.Run("changes the local part of an omni branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.False(t, bs.RemoteChanged())
		})
	})

	t.Run("RemoteRemoved", func(t *testing.T) {
		t.Parallel()
		t.Run("removing a remote-only branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.RemoteRemoved())
		})
		t.Run("removing the remote part of an omni branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.RemoteRemoved())
		})

		t.Run("changes a remote branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			must.False(t, bs.RemoteRemoved())
		})

		t.Run("upstream branch", func(t *testing.T) {
			t.Parallel()
			bs := undo.BranchSpan{
				Before: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("upstream/main"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			must.False(t, bs.RemoteRemoved())
		})
	})
}
