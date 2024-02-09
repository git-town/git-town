package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()

	t.Run("FindMatchingRecord", func(t *testing.T) {
		t.Parallel()
		t.Run("has matching local name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("branch-1"),
					LocalSHA:   gitdomain.NewSHA("111111"),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("branch-1"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			have := bis.FindMatchingRecord(give)
			want := bis[0]
			must.EqOp(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  gitdomain.NewSHA("111111"),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/branch-1"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
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
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("one"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.True(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("two"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with a matching remote", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("two"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: gitdomain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: gitdomain.NewRemoteBranchName("origin/one"),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("one"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.False(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("up-to-date"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("ahead"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: gitdomain.NewRemoteBranchName("origin/ahead"),
				RemoteSHA:  gitdomain.NewSHA("222222"),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("behind"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: gitdomain.NewRemoteBranchName("origin/behind"),
				RemoteSHA:  gitdomain.NewSHA("222222"),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("local-only"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/remote-only"),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		have := bs.LocalBranches().Names()
		want := gitdomain.NewLocalBranchNames("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("up-to-date"),
				LocalSHA:   gitdomain.NewSHA("111111"),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: gitdomain.NewRemoteBranchName("origin/up-to-date"),
				RemoteSHA:  gitdomain.NewSHA("111111"),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("ahead"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("behind"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("local-only"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("remote-only"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("deleted-at-remote"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().Names()
		want := gitdomain.NewLocalBranchNames("deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LookupLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch with matching name", func(t *testing.T) {
			branchOne := gitdomain.NewLocalBranchName("one")
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  branchOne,
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.EqOp(t, branchOne, bs.FindByLocalName(branchOne).LocalName)
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.EmptyLocalBranchName(),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.NewRemoteBranchName("kg/one"),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			have := bs.FindByLocalName(gitdomain.NewLocalBranchName("kg/one"))
			must.Nil(t, have)
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/two"),
				RemoteSHA:  gitdomain.EmptySHA(),
			}
			bs := gitdomain.BranchInfos{branch}
			have := bs.FindByRemoteName(gitdomain.NewRemoteBranchName("origin/two"))
			must.EqOp(t, branch, *have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("kg/one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			}}
			have := bs.FindByRemoteName(gitdomain.NewRemoteBranchName("kg/one"))
			must.Nil(t, have)
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("two"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.EmptyLocalBranchName(),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: gitdomain.NewRemoteBranchName("origin/three"),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		have := bs.Names()
		want := gitdomain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("one"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("two"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			have := bs.Remove(gitdomain.NewLocalBranchName("two"))
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchName("one"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("two"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		have := bs.Remove(gitdomain.NewLocalBranchName("zonk"))
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("two"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("two"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("three"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("four"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		have, err := bs.Select(gitdomain.LocalBranchNames{gitdomain.NewLocalBranchName("one"), gitdomain.NewLocalBranchName("three")})
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("one"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchName("three"),
				LocalSHA:   gitdomain.EmptySHA(),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: gitdomain.EmptyRemoteBranchName(),
				RemoteSHA:  gitdomain.EmptySHA(),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})
}
