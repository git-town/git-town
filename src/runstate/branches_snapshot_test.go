package runstate_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/stretchr/testify/assert"
)

func TestBranchBeforeAfter(t *testing.T) {
	t.Parallel()

	t.Run("IsLocalAdd", func(t *testing.T) {
		t.Parallel()
		t.Run("is a local add", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bba.IsLocalAdd())
		})
		t.Run("also adds a remote branch", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.False(t, bba.IsLocalAdd())
		})
		t.Run("doesn't add anything", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bba.IsLocalAdd())
		})
	})

	t.Run("IsLocalRemove", func(t *testing.T) {
		t.Parallel()
		t.Run("removes a local-only branch", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bba.IsLocalRemove())
		})
		t.Run("removes an omni branch", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bba.IsLocalRemove())
		})
		t.Run("changes a local-only branch", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.False(t, bba.IsLocalRemove())
		})
	})

	t.Run("IsOmniAdd", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omniadd", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.True(t, bba.IsOmniAdd())
		})
		t.Run("not an omniadd", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.False(t, bba.IsOmniAdd())
		})
	})

	t.Run("IsOmniChange", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omni change", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			assert.True(t, bba.IsOmniChange())
		})
		t.Run("not an omni change", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("333333"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			assert.False(t, bba.IsOmniChange())
		})
	})

	t.Run("IsOmniRemove", func(t *testing.T) {
		t.Parallel()
		t.Run("is an omniremove", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.LocalBranchName{},
					LocalSHA:   domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			assert.True(t, bba.IsOmniRemove())
		})
		t.Run("not an omniremove", func(t *testing.T) {
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.False(t, bba.IsOmniRemove())
		})
	})

	t.Run("LocalChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("has local changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.True(t, bba.LocalChanged())
		})
		t.Run("no local changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			assert.False(t, bba.LocalChanged())
		})
	})

	t.Run("NoChanges", func(t *testing.T) {
		t.Parallel()
		t.Run("no changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.True(t, bba.NoChanges())
		})
		t.Run("has changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			assert.False(t, bba.NoChanges())
		})
	})

	t.Run("RemoteChanged", func(t *testing.T) {
		t.Parallel()
		t.Run("has remote changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("222222"),
				},
			}
			assert.True(t, bba.RemoteChanged())
		})
		t.Run("no remote changes", func(t *testing.T) {
			t.Parallel()
			bba := runstate.BranchBeforeAfter{
				Before: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("222222"),
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
				After: domain.BranchInfo{
					LocalName:  domain.NewLocalBranchName("branch-1"),
					LocalSHA:   domain.NewSHA("111111"),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.NewSHA("111111"),
				},
			}
			assert.False(t, bba.RemoteChanged())
		})
	})

}

func TestBranchesSnapshot(t *testing.T) {
	t.Parallel()

	// t.Run("Diff", func(t *testing.T) {
	// 	t.Parallel()
	// 	t.Run("local-only branch added", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.RemoteBranchName{},
	// 					RemoteSHA:  domain.SHA{},
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:    domain.NewLocalBranchNames("branch-1"),
	// 			LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded:   []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("local-only branch removed", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.RemoteBranchName{},
	// 					RemoteSHA:  domain.SHA{},
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded: domain.LocalBranchNames{},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{
	// 				domain.NewLocalBranchName("branch-1"): domain.NewSHA("111111"),
	// 			},
	// 			LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded:   []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("local-only branch changed", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.RemoteBranchName{},
	// 					RemoteSHA:  domain.SHA{},
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("222222"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.RemoteBranchName{},
	// 					RemoteSHA:  domain.SHA{},
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:   domain.LocalBranchNames{},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{
	// 				domain.NewLocalBranchName("branch-1"): {
	// 					Before: domain.NewSHA("111111"),
	// 					After:  domain.NewSHA("222222"),
	// 				},
	// 			},
	// 			RemoteAdded:   []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("local-only branch pushed to origin", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.RemoteBranchName{},
	// 					RemoteSHA:  domain.SHA{},
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusLocalOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:   domain.LocalBranchNames{},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded: []domain.RemoteBranchName{
	// 				domain.NewRemoteBranchName("origin/branch-1"),
	// 			},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("remote-only branch added", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.LocalBranchName{},
	// 					LocalSHA:   domain.SHA{},
	// 					SyncStatus: domain.SyncStatusRemoteOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:   domain.LocalBranchNames{},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded: []domain.RemoteBranchName{
	// 				domain.NewRemoteBranchName("origin/branch-1"),
	// 			},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("remote-only branch downloaded", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.LocalBranchName{},
	// 					LocalSHA:   domain.SHA{},
	// 					SyncStatus: domain.SyncStatusRemoteOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusUpToDate,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded: domain.LocalBranchNames{
	// 				domain.NewLocalBranchName("branch-1"),
	// 			},
	// 			LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded:   []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("remote-only branch deleted", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.LocalBranchName{},
	// 					LocalSHA:   domain.SHA{},
	// 					SyncStatus: domain.SyncStatusRemoteOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:   domain.LocalBranchNames{},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded:  []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{
	// 				domain.NewRemoteBranchName("origin/branch-1"): domain.NewSHA("111111"),
	// 			},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("remote-only branch changed", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.LocalBranchName{},
	// 					LocalSHA:   domain.SHA{},
	// 					SyncStatus: domain.SyncStatusRemoteOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.LocalBranchName{},
	// 					LocalSHA:   domain.SHA{},
	// 					SyncStatus: domain.SyncStatusRemoteOnly,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("222222"),
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded:    domain.LocalBranchNames{},
	// 			LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded:   []domain.RemoteBranchName{},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{
	// 				domain.NewRemoteBranchName("origin/branch-1"): {
	// 					Before: domain.NewSHA("111111"),
	// 					After:  domain.NewSHA("222222"),
	// 				},
	// 			},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	// 	t.Run("omnibranch added", func(t *testing.T) {
	// 		t.Parallel()
	// 		before := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{},
	// 		}
	// 		after := runstate.BranchesSnapshot{
	// 			Branches: domain.BranchInfos{
	// 				domain.BranchInfo{
	// 					LocalName:  domain.NewLocalBranchName("branch-1"),
	// 					LocalSHA:   domain.NewSHA("111111"),
	// 					SyncStatus: domain.SyncStatusUpToDate,
	// 					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
	// 					RemoteSHA:  domain.NewSHA("111111"),
	// 				},
	// 			},
	// 		}
	// 		have := before.Diff(after)
	// 		want := runstate.Changes{
	// 			LocalAdded: domain.LocalBranchNames{
	// 				domain.NewLocalBranchName("branch-1"),
	// 			},
	// 			LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
	// 			LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
	// 			RemoteAdded: []domain.RemoteBranchName{
	// 				domain.NewRemoteBranchName("origin/branch-1"),
	// 			},
	// 			RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
	// 			RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
	// 		}
	// 		assert.Equal(t, want, have)
	// 	})

	t.Run("omnibranch changed locally", func(t *testing.T) {})
	t.Run("omnibranch changed remotely", func(t *testing.T) {})
	t.Run("omnibranch changed locally and remotely to same SHA", func(t *testing.T) {})
	t.Run("omnibranch changed locally and remotely to different SHAs", func(t *testing.T) {})
	t.Run("omnibranch updates pulled down", func(t *testing.T) {})
	t.Run("omnibranch updates pushed up", func(t *testing.T) {})
	t.Run("omnibranch deleted locally", func(t *testing.T) {})
	t.Run("omnibranch deleted remotely", func(t *testing.T) {})
	// })
}
