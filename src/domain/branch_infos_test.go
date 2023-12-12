package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()

	t.Run("FindMatchingRecord", func(t *testing.T) {
		t.Parallel()
		t.Run("has matching local name", func(t *testing.T) {
			t.Parallel()
			bis := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			give := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}
			have := bis.FindMatchingRecord(give)
			want := bis[0]
			must.EqOp(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			give := domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  domain.NewSHA("111111"),
			}
			have := bis.FindMatchingRecord(give)
			want := bis[0]
			must.EqOp(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.HasLocalBranch(domain.NewLocalBranchName("one")))
		})
	})

	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with a matching remote", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(domain.NewLocalBranchName("one")))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusRemoteOnly,
					RemoteName: domain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(domain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.False(t, bs.HasMatchingTrackingBranchFor(domain.NewLocalBranchName("one")))
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
				SyncStatus: domain.SyncStatusNotInSync,
				RemoteName: domain.NewRemoteBranchName("origin/ahead"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("behind"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusNotInSync,
				RemoteName: domain.NewRemoteBranchName("origin/behind"),
				RemoteSHA:  domain.NewSHA("222222"),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("local-only"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/remote-only"),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusDeletedAtRemote,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		have := bs.LocalBranches().Names()
		want := domain.NewLocalBranchNames("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		must.Eq(t, want, have)
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
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusNotInSync,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("behind"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusNotInSync,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("local-only"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("remote-only"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusDeletedAtRemote,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().Names()
		want := domain.NewLocalBranchNames("deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LookupLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch with matching name", func(t *testing.T) {
			branchOne := domain.NewLocalBranchName("one")
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  branchOne,
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.EqOp(t, branchOne, bs.FindByLocalName(branchOne).LocalName)
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.EmptyLocalBranchName(),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("kg/one"),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			have := bs.FindByLocalName(domain.NewLocalBranchName("kg/one"))
			must.Nil(t, have)
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.NewRemoteBranchName("origin/two"),
				RemoteSHA:  domain.EmptySHA(),
			}
			bs := domain.BranchInfos{branch}
			have := bs.FindByRemoteName(domain.NewRemoteBranchName("origin/two"))
			must.EqOp(t, branch, *have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := domain.BranchInfos{domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("kg/one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			}}
			have := bs.FindByRemoteName(domain.NewRemoteBranchName("kg/one"))
			must.Nil(t, have)
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.EmptyLocalBranchName(),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusRemoteOnly,
				RemoteName: domain.NewRemoteBranchName("origin/three"),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		have := bs.Names()
		want := domain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("two"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			have := bs.Remove(domain.NewLocalBranchName("two"))
			want := domain.BranchInfos{
				domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("one"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		have := bs.Remove(domain.NewLocalBranchName("zonk"))
		want := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("two"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("three"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("four"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		have, err := bs.Select(domain.LocalBranchNames{domain.NewLocalBranchName("one"), domain.NewLocalBranchName("three")})
		want := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("one"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("three"),
				LocalSHA:   domain.EmptySHA(),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.EmptyRemoteBranchName(),
				RemoteSHA:  domain.EmptySHA(),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})
}
