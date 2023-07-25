package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestBranches(t *testing.T) {
	t.Parallel()

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusAhead,
			},
		}
		have := bs.BranchNames()
		want := []string{"one", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusAhead,
			},
		}
		assert.True(t, bs.Contains("one"))
		assert.True(t, bs.Contains("two"))
		assert.False(t, bs.Contains("zonk"))
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.BranchSyncStatus{
				Name:       "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.BranchSyncStatus{
				Name:       "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.BranchSyncStatus{
				Name:       "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.BranchSyncStatus{
				Name:       "deleted-at-remote",
				SyncStatus: git.SyncStatusDeletedAtRemote,
			},
		}
		want := []string{"up-to-date", "ahead", "behind", "local-only", "deleted-at-remote"}
		have := bs.LocalBranches().BranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.BranchSyncStatus{
				Name:       "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.BranchSyncStatus{
				Name:       "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.BranchSyncStatus{
				Name:       "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.BranchSyncStatus{
				Name:       "deleted-at-remote",
				SyncStatus: git.SyncStatusDeletedAtRemote,
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().BranchNames()
		want := []string{"deleted-at-remote"}
		assert.Equal(t, want, have)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.Equal(t, "one", bs.Lookup("one").Name)
		assert.Equal(t, "two", bs.Lookup("two").Name)
		assert.Nil(t, bs.Lookup("zonk"))
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:       "one",
					SyncStatus: git.SyncStatusUpToDate,
				},
				git.BranchSyncStatus{
					Name:       "two",
					SyncStatus: git.SyncStatusUpToDate,
				},
			}
			have := bs.Remove("two")
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:       "one",
					SyncStatus: git.SyncStatusUpToDate,
				},
			}
			assert.Equal(t, want, have)
		})
	})
	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		have := bs.Remove("zonk")
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "three",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "four",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		give := []string{"one", "three"}
		have, err := bs.Select(give)
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:       "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				Name:       "three",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
