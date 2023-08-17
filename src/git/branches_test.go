package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestBranch(t *testing.T) {
	t.Parallel()
	t.Run("RemoteBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("remote-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("origin/branch1"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			}
			have := give.RemoteBranch()
			want := domain.NewRemoteBranchName("origin/branch1")
			assert.Equal(t, want, have)
		})
		t.Run("local-only branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("branch1"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			}
			have := give.RemoteBranch()
			want := domain.RemoteBranchName{}
			assert.Equal(t, want, have)
		})
		t.Run("local branch with tracking branch", func(t *testing.T) {
			t.Parallel()
			give := git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("branch1"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: domain.NewRemoteBranchName("origin/branch-2"),
				TrackingSHA:  domain.SHA{},
			}
			have := give.RemoteBranch()
			want := domain.NewRemoteBranchName("origin/branch-2")
			assert.Equal(t, want, have)
		})
	})
}

func TestBranches(t *testing.T) {
	t.Parallel()

	// t.Run("BranchNames", func(t *testing.T) {
	// 	t.Parallel()
	// 	bs := git.BranchesSyncStatus{
	// 		git.BranchSyncStatus{
	// 			Name:         domain.NewLocalBranchName("one"),
	// 			InitialSHA:   domain.SHA{},
	// 			SyncStatus:   git.SyncStatusLocalOnly,
	// 			TrackingName: domain.RemoteBranchName{},
	// 			TrackingSHA:  domain.SHA{},
	// 		},
	// 		git.BranchSyncStatus{
	// 			Name:         domain.NewLocalBranchName("two"),
	// 			InitialSHA:   domain.SHA{},
	// 			SyncStatus:   git.SyncStatusLocalOnly,
	// 			TrackingName: domain.RemoteBranchName{},
	// 			TrackingSHA:  domain.SHA{},
	// 		},
	// 	}
	// 	have := bs.Branch
	// 	want := []string{"one", "two"}
	// 	assert.Equal(t, want, have)
	// })

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("one"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.RemoteBranchName{},
					TrackingSHA:  domain.SHA{},
				},
			}
			assert.True(t, bs.ContainsLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("origin/one"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.RemoteBranchName{},
					TrackingSHA:  domain.SHA{},
				},
			}
			assert.False(t, bs.ContainsLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("two"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.NewRemoteBranchName("origin/one"),
					TrackingSHA:  domain.SHA{},
				},
			}
			assert.False(t, bs.ContainsLocalBranch(domain.NewLocalBranchName("one")))
		})
	})

	// t.Run("IsKnown", func(t *testing.T) {
	// 	t.Parallel()
	// 	t.Run("the branch in question is a local branch", func(t *testing.T) {
	// 		bs := git.BranchesSyncStatus{
	// 			git.BranchSyncStatus{
	// 				Name:         domain.NewLocalBranchName("one"),
	// 				InitialSHA:   domain.SHA{},
	// 				SyncStatus:   git.SyncStatusLocalOnly,
	// 				TrackingName: domain.RemoteBranchName{},
	// 				TrackingSHA:  domain.SHA{},
	// 			},
	// 			git.BranchSyncStatus{
	// 				Name:         domain.NewLocalBranchName("two"),
	// 				InitialSHA:   domain.SHA{},
	// 				SyncStatus:   git.SyncStatusLocalOnly,
	// 				TrackingName: domain.RemoteBranchName{},
	// 				TrackingSHA:  domain.SHA{},
	// 			},
	// 		}
	// 		assert.True(t, bs.ContainsLocalBranch("one"))
	// 		assert.True(t, bs.IsKnown("two"))
	// 		assert.False(t, bs.IsKnown("zonk"))
	// 	})
	// 	t.Run("the branch in question is a tracking branch of an already known local branch", func(t *testing.T) {
	// 		bs := git.BranchesSyncStatus{
	// 			git.BranchSyncStatus{
	// 				Name:         domain.NewLocalBranchName("one"),
	// 				InitialSHA:   domain.SHA{},
	// 				SyncStatus:   git.SyncStatusUpToDate,
	// 				TrackingName: domain.NewRemoteBranchName("origin/two"),
	// 				TrackingSHA:  domain.SHA{},
	// 			},
	// 		}
	// 		assert.True(t, bs.IsKnown("origin/two"))
	// 		assert.False(t, bs.IsKnown("zonk"))
	// 	})
	// })

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("up-to-date"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: domain.NewRemoteBranchName("origin/up-to-date"),
				TrackingSHA:  domain.NewSHA("111111"),
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("ahead"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: domain.NewRemoteBranchName("origin/ahead"),
				TrackingSHA:  domain.NewSHA("222222"),
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("behind"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: domain.NewRemoteBranchName("origin/behind"),
				TrackingSHA:  domain.NewSHA("222222"),
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("local-only"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("remote-only"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("deleted-at-remote"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		want := domain.LocalBranchNamesFrom("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		have := bs.LocalBranches().LocalBranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("up-to-date"),
				InitialSHA:   domain.NewSHA("111111"),
				SyncStatus:   git.SyncStatusUpToDate,
				TrackingName: domain.NewRemoteBranchName("origin/up-to-date"),
				TrackingSHA:  domain.NewSHA("111111"),
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("ahead"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusAhead,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("behind"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusBehind,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("local-only"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("remote-only"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusRemoteOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("deleted-at-remote"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusDeletedAtRemote,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().LocalBranchNames()
		want := domain.LocalBranchNamesFrom("deleted-at-remote")
		assert.Equal(t, want, have)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		branchOne := domain.NewLocalBranchName("one")
		branchTwo := domain.NewLocalBranchName("two")
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         branchOne,
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         branchTwo,
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		assert.Equal(t, branchOne, bs.Lookup(branchOne).Name)
		assert.Equal(t, branchTwo, bs.Lookup(branchTwo).Name)
		assert.Nil(t, bs.Lookup(domain.NewLocalBranchName("zonk")))
	})

	t.Run("LookupLocalBranchWithTracking", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.NewRemoteBranchName("origin/two"),
				TrackingSHA:  domain.SHA{},
			}
			bs := git.BranchesSyncStatus{branch}
			have := bs.LookupLocalBranchWithTracking(domain.NewRemoteBranchName("origin/two"))
			assert.Equal(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := git.BranchesSyncStatus{git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("kg/one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			}}
			have := bs.LookupLocalBranchWithTracking(domain.NewRemoteBranchName("kg/one"))
			assert.Nil(t, have)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("one"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.RemoteBranchName{},
					TrackingSHA:  domain.SHA{},
				},
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("two"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.RemoteBranchName{},
					TrackingSHA:  domain.SHA{},
				},
			}
			have := bs.Remove(domain.NewLocalBranchName("two"))
			want := git.BranchesSyncStatus{
				git.BranchSyncStatus{
					Name:         domain.NewLocalBranchName("one"),
					InitialSHA:   domain.SHA{},
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: domain.RemoteBranchName{},
					TrackingSHA:  domain.SHA{},
				},
			}
			assert.Equal(t, want, have)
		})
	})
	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("two"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		have := bs.Remove(domain.NewLocalBranchName("zonk"))
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("two"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("two"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("three"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("four"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		have, err := bs.Select([]domain.LocalBranchName{domain.NewLocalBranchName("one"), domain.NewLocalBranchName("three")})
		want := git.BranchesSyncStatus{
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("one"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
			git.BranchSyncStatus{
				Name:         domain.NewLocalBranchName("three"),
				InitialSHA:   domain.SHA{},
				SyncStatus:   git.SyncStatusLocalOnly,
				TrackingName: domain.RemoteBranchName{},
				TrackingSHA:  domain.SHA{},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
