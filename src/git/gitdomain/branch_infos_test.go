package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
					LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("branch-1")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			have := bis.FindMatchingRecord(give)
			want := Some(bis[0])
			must.Eq(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			have := bis.FindMatchingRecord(give)
			want := Some(bis[0])
			must.Eq(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("one")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.True(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("two")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
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
					LocalName:  Some(gitdomain.NewLocalBranchName("two")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("one")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.False(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("up-to-date")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("ahead")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("behind")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("local-only")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("deleted-at-remote")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
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
				LocalName:  Some(gitdomain.NewLocalBranchName("up-to-date")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("ahead")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("behind")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("local-only")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("deleted-at-remote")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
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
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(branchOne),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			branchInfo, hasBranchInfo := branchInfos.FindByLocalName(branchOne).Get()
			must.True(t, hasBranchInfo)
			must.EqOp(t, branchOne, branchInfo.LocalName.GetOrPanic())
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("kg/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			have := branchInfos.FindByLocalName(gitdomain.NewLocalBranchName("kg/one"))
			must.True(t, have.IsNone())
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("two")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/two")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bs := gitdomain.BranchInfos{branch}
			have, has := bs.FindByRemoteName(gitdomain.NewRemoteBranchName("origin/two")).Get()
			must.True(t, has)
			must.Eq(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("kg/one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}}
			have := bs.FindByRemoteName(gitdomain.NewRemoteBranchName("kg/one"))
			must.True(t, have.IsNone())
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("two")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/three")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
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
					LocalName:  Some(gitdomain.NewLocalBranchName("one")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("two")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			have := bs.Remove(gitdomain.NewLocalBranchName("two"))
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(gitdomain.NewLocalBranchName("one")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("two")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.Remove(gitdomain.NewLocalBranchName("zonk"))
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("two")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("two")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("three")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("four")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have, err := bs.Select(gitdomain.NewLocalBranchName("one"), gitdomain.NewLocalBranchName("three"))
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("one")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  Some(gitdomain.NewLocalBranchName("three")),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})
}
