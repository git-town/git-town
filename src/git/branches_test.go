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
				Name:         "branch1",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			}
			have := give.NameWithoutRemote()
			want := "branch1"
			assert.Equal(t, want, have)
		})
		t.Run("remote branch", func(t *testing.T) {
			give := git.BranchSyncStatus{
				Name:         "origin/branch1",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  "",
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
				Name:         "origin/branch1",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  "",
			}
			have := give.RemoteBranch()
			want := "origin/branch1"
			assert.Equal(t, want, have)
		})
		t.Run("local-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         "branch1",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			}
			have := give.RemoteBranch()
			want := ""
			assert.Equal(t, want, have)
		})
		t.Run("local branch with tracking branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         "branch1",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/branch-2",
				TrackingSHA:  "",
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
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
		}
		have := bs.BranchNames()
		want := []string{"one", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("IsKnown", func(t *testing.T) {
		t.Parallel()
		t.Run("the branch in question is a local branch", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
				git.BranchSyncStatus{
					Name:         "two",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
			}
			assert.True(t, bs.IsKnown("one"))
			assert.True(t, bs.IsKnown("two"))
			assert.False(t, bs.IsKnown("zonk"))
		})
		t.Run("the branch in question is a tracking branch of an already known local branch", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusUpToDate,
					TrackingName: "origin/two",
					TrackingSHA:  "",
				},
			}
			assert.True(t, bs.IsKnown("origin/two"))
			assert.False(t, bs.IsKnown("zonk"))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "up-to-date",
				InitialSHA:   "11111111",
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/up-to-date",
				TrackingSHA:  "11111111",
			},
			git.BranchSyncStatus{
				Name:         "ahead",
				InitialSHA:   "11111111",
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: "origin/ahead",
				TrackingSHA:  "22222222",
			},
			git.BranchSyncStatus{
				Name:         "behind",
				InitialSHA:   "111111111",
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: "origin/behind",
				TrackingSHA:  "222222222",
			},
			git.BranchSyncStatus{
				Name:         "local-only",
				InitialSHA:   "11111111",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "remote-only",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "deleted-at-remote",
				InitialSHA:   "11111111111",
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: "",
				TrackingSHA:  "",
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
				Name:         "up-to-date",
				InitialSHA:   "1111111111",
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/up-to-date",
				TrackingSHA:  "1111111111",
			},
			git.BranchSyncStatus{
				Name:         "ahead",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: "origin/ahead",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "behind",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: "origin/behind",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "local-only",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "remote-only",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "deleted-at-remote",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: "",
				TrackingSHA:  "",
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
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
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
					Name:         "one",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
				git.BranchSyncStatus{
					Name:         "two",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
			}
			have := bs.Remove("two")
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
			}
			assert.Equal(t, want, have)
		})
	})
	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
		}
		have := bs.Remove("zonk")
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "three",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "four",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
		}
		have, err := bs.Select([]string{"one", "three"})
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
			git.BranchSyncStatus{
				Name:         "three",
				InitialSHA:   "",
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  "",
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
