package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfo(t *testing.T) {
	t.Parallel()

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("one"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.LocalBranchName{},
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("two"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
	})

	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with a matching remote", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("two"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bs.HasMatchingRemoteBranchFor(domain.NewLocalBranchName("one")))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.LocalBranchName{},
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bs.HasMatchingRemoteBranchFor(domain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("one"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bs.HasMatchingRemoteBranchFor(domain.NewLocalBranchName("one")))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("up-to-date"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  domain.NewSHA("111111"),
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("ahead"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusAhead,
				RemoteName: domain.NewRemoteBranchName("origin/ahead"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("behind"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusBehind,
				RemoteName: domain.NewRemoteBranchName("origin/behind"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("local-only"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("remote-only"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("deleted-at-remote"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusDeletedAtRemote,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		want := domain.NewLocalBranchNames("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		have := bs.LocalBranches().Names()
		assert.Equal(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("up-to-date"),
				InitialSHA: domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  domain.NewSHA("111111"),
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("ahead"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusAhead,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("behind"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusBehind,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("local-only"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("remote-only"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("deleted-at-remote"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusDeletedAtRemote,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().Names()
		want := domain.NewLocalBranchNames("deleted-at-remote")
		assert.Equal(t, want, have)
	})

	t.Run("LookupLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch with matching name", func(t *testing.T) {
			branchOne := domain.NewLocalBranchName("one")
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       branchOne,
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.Equal(t, branchOne, bs.FindLocalBranch(branchOne).Name)
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.LocalBranchName{},
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("kg/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			have := bs.FindLocalBranch(domain.NewLocalBranchName("kg/one"))
			assert.Nil(t, have)
		})
	})

	t.Run("LookupLocalBranchWithTracking", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/two"),
				RemoteSHA:  domain.SHA{},
			}
			bs := domain.BranchInfos{branch}
			have := bs.FindLocalBranchWithTracking(domain.NewRemoteBranchName("origin/two"))
			assert.Equal(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{domain.BranchInfo{
				Name:       domain.NewLocalBranchName("kg/one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			}}
			have := bs.FindLocalBranchWithTracking(domain.NewRemoteBranchName("kg/one"))
			assert.Nil(t, have)
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("two"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.LocalBranchName{},
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/three"),
				RemoteSHA:  domain.SHA{},
			},
		}
		have := bs.Names()
		want := domain.NewLocalBranchNames("one", "two")
		assert.Equal(t, want, have)
	})

	t.Run("Remote", func(t *testing.T) {
		t.Parallel()
		t.Run("Remote branch is set", func(t *testing.T) {
			branchInfo := domain.BranchInfo{
				Name:       domain.LocalBranchName{},
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch"),
				RemoteSHA:  domain.SHA{},
			}
			have := branchInfo.Remote()
			want := "origin"
			assert.Equal(t, want, have)
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("one"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("two"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			have := bs.Remove(domain.NewLocalBranchName("two"))
			want := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("one"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.Equal(t, want, have)
		})
	})
	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("two"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		have := bs.Remove(domain.NewLocalBranchName("zonk"))
		want := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("two"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		assert.Equal(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("two"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("three"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("four"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		have, err := bs.Select([]domain.LocalBranchName{domain.NewLocalBranchName("one"), domain.NewLocalBranchName("three")})
		want := domain.BranchInfos{
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("one"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				Name:       domain.NewLocalBranchName("three"),
				InitialSHA: domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
