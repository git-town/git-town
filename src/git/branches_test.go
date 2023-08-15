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
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			}
			have := give.NameWithoutRemote()
			want := "branch1"
			assert.Equal(t, want, have)
		})
		t.Run("remote branch", func(t *testing.T) {
			give := git.BranchSyncStatus{
				Name:         "origin/branch1",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			}
			have := give.RemoteBranch()
			want := "origin/branch1"
			assert.Equal(t, want, have)
		})
		t.Run("local-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         "branch1",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			}
			have := give.RemoteBranch()
			want := ""
			assert.Equal(t, want, have)
		})
		t.Run("local branch with tracking branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         "branch1",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/branch-2",
				TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		have := bs.BranchNames()
		want := []string{"one", "two"}
		assert.Equal(t, want, have)
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
				},
			}
			assert.True(t, bs.HasLocalBranch("one"))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "origin/one",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
				},
			}
			assert.False(t, bs.HasLocalBranch("one"))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "two",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "origin/one",
					TrackingSHA:  git.SHA{""},
				},
			}
			assert.False(t, bs.HasLocalBranch("one"))
		})
	})

	t.Run("IsKnown", func(t *testing.T) {
		t.Parallel()
		t.Run("the branch in question is a local branch", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
				},
				git.BranchSyncStatus{
					Name:         "two",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
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
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusUpToDate,
					TrackingName: "origin/two",
					TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{"11111111"},
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/up-to-date",
				TrackingSHA:  git.SHA{"11111111"},
			},
			git.BranchSyncStatus{
				Name:         "ahead",
				InitialSHA:   git.SHA{"11111111"},
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: "origin/ahead",
				TrackingSHA:  git.SHA{"22222222"},
			},
			git.BranchSyncStatus{
				Name:         "behind",
				InitialSHA:   git.SHA{"111111111"},
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: "origin/behind",
				TrackingSHA:  git.SHA{"222222222"},
			},
			git.BranchSyncStatus{
				Name:         "local-only",
				InitialSHA:   git.SHA{"11111111"},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "remote-only",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "deleted-at-remote",
				InitialSHA:   git.SHA{"11111111111"},
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{"1111111111"},
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: "origin/up-to-date",
				TrackingSHA:  git.SHA{"1111111111"},
			},
			git.BranchSyncStatus{
				Name:         "ahead",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: "origin/ahead",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "behind",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: "origin/behind",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "local-only",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "remote-only",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "deleted-at-remote",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		assert.Equal(t, "one", bs.Lookup("one").Name)
		assert.Equal(t, "two", bs.Lookup("two").Name)
		assert.Nil(t, bs.Lookup("zonk"))
	})

	t.Run("LookupLocalBranchWithTracking", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			branch := git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "origin/two",
				TrackingSHA:  git.SHA{""},
			}
			bs := git.BranchesSyncStatus{branch}
			have := bs.LookupLocalBranchWithTracking("origin/two")
			assert.Equal(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			bs := git.BranchesSyncStatus{git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			}}
			have := bs.LookupLocalBranchWithTracking("one")
			assert.Nil(t, have)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
				},
				git.BranchSyncStatus{
					Name:         "two",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
				},
			}
			have := bs.Remove("two")
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         "one",
					InitialSHA:   git.SHA{""},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  git.SHA{""},
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
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		have := bs.Remove("zonk")
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "two",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "three",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "four",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		have, err := bs.Select([]string{"one", "three"})
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         "one",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
			git.BranchSyncStatus{
				Name:         "three",
				InitialSHA:   git.SHA{""},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: "",
				TrackingSHA:  git.SHA{""},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
