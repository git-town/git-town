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
	t.Run("NameWithoutRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch", func(t *testing.T) {
			give := git.BranchSyncStatus{
				Name:           "branch1",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			}
			have := give.NameWithoutRemote()
			want := "branch1"
			assert.Equal(t, want, have)
		})
		t.Run("remote branch", func(t *testing.T) {
			give := git.BranchSyncStatus{
				Name:           "origin/branch1",
				SyncStatus:     git.SyncStatusRemoteOnly,
				TrackingBranch: "",
			}
			have := give.NameWithoutRemote()
			want := "branch1"
			assert.Equal(t, want, have)
		})
	})
	t.Run("RemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("remote-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:           "origin/branch1",
				SyncStatus:     git.SyncStatusRemoteOnly,
				TrackingBranch: "",
			}
			have := give.RemoteBranch()
			want := "origin/branch1"
			assert.Equal(t, want, have)
		})
		t.Run("local-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:           "branch1",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			}
			have := give.RemoteBranch()
			want := ""
			assert.Equal(t, want, have)
		})
		t.Run("local branch with tracking branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:           "branch1",
				SyncStatus:     git.SyncStatusUpToDate,
				TrackingBranch: "origin/branch-2",
			}
			have := give.RemoteBranch()
			want := "origin/branch-2"
			assert.Equal(t, want, have)
		})
	})
}

func TestBranches(t *testing.T) {
	t.Parallel()

	t.Run("BranchNames", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "two",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
		}
		have := bs.BranchNames()
		want := []string{"one", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the branch directly", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:           "one",
					SyncStatus:     git.SyncStatusLocalOnly,
					TrackingBranch: "",
				},
				git.BranchSyncStatus{
					Name:           "two",
					SyncStatus:     git.SyncStatusLocalOnly,
					TrackingBranch: "",
				},
			}
			assert.True(t, bs.Contains("one"))
			assert.True(t, bs.Contains("two"))
			assert.False(t, bs.Contains("zonk"))
		})
		t.Run("contains a branch that has this branch as the tracking branch", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:           "one",
					SyncStatus:     git.SyncStatusUpToDate,
					TrackingBranch: "origin/two",
				},
			}
			assert.True(t, bs.Contains("origin/two"))
			assert.False(t, bs.Contains("zonk"))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "up-to-date",
				SyncStatus:     git.SyncStatusUpToDate,
				TrackingBranch: "origin/up-to-date",
			},
			git.BranchSyncStatus{
				Name:           "ahead",
				SyncStatus:     git.SyncStatusAhead,
				TrackingBranch: "origin/ahead",
			},
			git.BranchSyncStatus{
				Name:           "behind",
				SyncStatus:     git.SyncStatusBehind,
				TrackingBranch: "origin/behind",
			},
			git.BranchSyncStatus{
				Name:           "local-only",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "remote-only",
				SyncStatus:     git.SyncStatusRemoteOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "deleted-at-remote",
				SyncStatus:     git.SyncStatusDeletedAtRemote,
				TrackingBranch: "",
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
				Name:           "up-to-date",
				SyncStatus:     git.SyncStatusUpToDate,
				TrackingBranch: "origin/up-to-date",
			},
			git.BranchSyncStatus{
				Name:           "ahead",
				SyncStatus:     git.SyncStatusAhead,
				TrackingBranch: "origin/ahead",
			},
			git.BranchSyncStatus{
				Name:           "behind",
				SyncStatus:     git.SyncStatusBehind,
				TrackingBranch: "origin/behind",
			},
			git.BranchSyncStatus{
				Name:           "local-only",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "remote-only",
				SyncStatus:     git.SyncStatusRemoteOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "deleted-at-remote",
				SyncStatus:     git.SyncStatusDeletedAtRemote,
				TrackingBranch: "",
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
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "two",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
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
					Name:           "one",
					SyncStatus:     git.SyncStatusLocalOnly,
					TrackingBranch: "",
				},
				git.BranchSyncStatus{
					Name:           "two",
					SyncStatus:     git.SyncStatusLocalOnly,
					TrackingBranch: "",
				},
			}
			have := bs.Remove("two")
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:           "one",
					SyncStatus:     git.SyncStatusLocalOnly,
					TrackingBranch: "",
				},
			}
			assert.Equal(t, want, have)
		})
	})
	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "two",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
		}
		have := bs.Remove("zonk")
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "two",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "two",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "three",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "four",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
		}
		have, err := bs.Select([]string{"one", "three"})
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:           "one",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
			git.BranchSyncStatus{
				Name:           "three",
				SyncStatus:     git.SyncStatusLocalOnly,
				TrackingBranch: "",
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
