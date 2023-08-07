package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestTrackingBranchName(t *testing.T) {
	t.Parallel()
	give := "branch1"
	have := git.TrackingBranchName(give)
	want := "origin/branch1"
	assert.Equal(t, want, have)
}

func TestBranch(t *testing.T) {
	t.Parallel()
	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		give := git.BranchSyncStatus{
			LocalName:  "branch1",
			SyncStatus: git.SyncStatusUpToDate,
		}
		have := give.TrackingBranch()
		want := "origin/branch1"
		assert.Equal(t, want, have)
	})
}

func TestBranches(t *testing.T) {
	t.Parallel()

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
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
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
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
				LocalName:  "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.BranchSyncStatus{
				LocalName:  "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.BranchSyncStatus{
				LocalName:  "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.BranchSyncStatus{
				LocalName:  "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.BranchSyncStatus{
				LocalName:  "deleted-at-remote",
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
				LocalName:  "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.BranchSyncStatus{
				LocalName:  "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.BranchSyncStatus{
				LocalName:  "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.BranchSyncStatus{
				LocalName:  "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.BranchSyncStatus{
				LocalName:  "deleted-at-remote",
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
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.Equal(t, "one", bs.Lookup("one").LocalName)
		assert.Equal(t, "two", bs.Lookup("two").LocalName)
		assert.Nil(t, bs.Lookup("zonk"))
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					LocalName:  "one",
					SyncStatus: git.SyncStatusUpToDate,
				},
				git.BranchSyncStatus{
					LocalName:  "two",
					SyncStatus: git.SyncStatusUpToDate,
				},
			}
			have := bs.Remove("two")
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					LocalName:  "one",
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
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		have := bs.Remove("zonk")
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "two",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "three",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "four",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		give := []string{"one", "three"}
		have, err := bs.Select(give)
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				LocalName:  "one",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchSyncStatus{
				LocalName:  "three",
				SyncStatus: git.SyncStatusUpToDate,
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
