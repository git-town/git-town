package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()
	t.Run("Clone", func(t *testing.T) {
		t.Parallel()
		bis := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			},
		}
		have := bis.Clone()
		have[0].LocalName = domain.NewLocalBranchName("branch-2")
		have[0].LocalSHA = domain.NewSHA("222222")
		have[0].RemoteName = domain.NewRemoteBranchName("origin/branch-2")
		have[0].RemoteSHA = domain.NewSHA("222222")
		assert.Equal(t, bis[0].LocalName, domain.NewLocalBranchName("branch-1"))
		assert.Equal(t, bis[0].LocalSHA, domain.NewSHA("111111"))
		assert.Equal(t, bis[0].RemoteName, domain.NewRemoteBranchName("origin/branch-1"))
		assert.Equal(t, bis[0].RemoteSHA, domain.NewSHA("111111"))
	})

	t.Run("FindMatchingRecord", func(t *testing.T) {
		t.Parallel()
		t.Run("has matching local name", func(t *testing.T) {
			t.Parallel()
			bis := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			}
			have := bis.FindMatchingRecord(give)
			want := bis[0]
			assert.Equal(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			give := domain.BranchInfo{
				LocalName:  domain.LocalBranchName{},
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			have := bis.FindMatchingRecord(give)
			want := bis[0]
			assert.Equal(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.SHA{},
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
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
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
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.SHA{},
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
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.SHA{},
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
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
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
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.SHA{},
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
				LocalName:  domain.NewLocalBranchName("up-to-date"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  domain.NewSHA("111111"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("ahead"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusAhead,
				RemoteName: domain.NewRemoteBranchName("origin/ahead"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("behind"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusBehind,
				RemoteName: domain.NewRemoteBranchName("origin/behind"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("local-only"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("remote-only"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   domain.NewSHA("111111"),
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
				LocalName:  domain.NewLocalBranchName("up-to-date"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusUpToDate,
				RemoteName: domain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  domain.NewSHA("111111"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("ahead"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusAhead,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("behind"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusBehind,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("local-only"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("remote-only"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   domain.SHA{},
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
					LocalName:  branchOne,
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.Equal(t, branchOne, bs.FindLocalBranch(branchOne).LocalName)
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("kg/one"),
					RemoteSHA:  domain.SHA{},
				},
			}
			have := bs.FindLocalBranch(domain.NewLocalBranchName("kg/one"))
			assert.Nil(t, have)
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/two"),
				RemoteSHA:  domain.SHA{},
			}
			bs := domain.BranchInfos{branch}
			have := bs.FindByRemote(domain.NewRemoteBranchName("origin/two"))
			assert.Equal(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("kg/one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			}}
			have := bs.FindByRemote(domain.NewRemoteBranchName("kg/one"))
			assert.Nil(t, have)
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.LocalBranchName{},
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/three"),
				RemoteSHA:  domain.SHA{},
			},
		}
		have := bs.Names()
		want := domain.NewLocalBranchNames("one", "two")
		assert.Equal(t, want, have)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			have := bs.Remove(domain.NewLocalBranchName("two"))
			want := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.SHA{},
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
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		have := bs.Remove(domain.NewLocalBranchName("zonk"))
		want := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.SHA{},
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
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("three"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("four"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		have, err := bs.Select([]domain.LocalBranchName{domain.NewLocalBranchName("one"), domain.NewLocalBranchName("three")})
		want := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("three"),
				LocalSHA:   domain.SHA{},
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, have, want)
	})
}
