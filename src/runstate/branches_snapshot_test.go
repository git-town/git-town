package runstate_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/stretchr/testify/assert"
)

func TestBranchesSnapshot(t *testing.T) {
	t.Parallel()
	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("local-only branch added", func(t *testing.T) {
			t.Parallel()
			before := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{},
			}
			after := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			have := before.Diff(after)
			want := runstate.BranchesDiff{
				LocalAdded:    domain.NewLocalBranchNames("branch-1"),
				LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
				LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local-only branch removed", func(t *testing.T) {
			t.Parallel()
			before := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			after := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{},
			}
			have := before.Diff(after)
			want := runstate.BranchesDiff{
				LocalAdded: domain.LocalBranchNames{},
				LocalRemoved: map[domain.LocalBranchName]domain.SHA{
					domain.NewLocalBranchName("branch-1"): domain.NewSHA("111111"),
				},
				LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local-only branch changed", func(t *testing.T) {
			t.Parallel()
			before := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			after := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("222222"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			have := before.Diff(after)
			want := runstate.BranchesDiff{
				LocalAdded:   domain.LocalBranchNames{},
				LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
				LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{
					domain.NewLocalBranchName("branch-1"): {
						Before: domain.NewSHA("111111"),
						After:  domain.NewSHA("222222"),
					},
				},
				RemoteAdded:   []domain.RemoteBranchName{},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			}
			assert.Equal(t, want, have)
		})

		t.Run("local-only branch pushed to origin", func(t *testing.T) {
			t.Parallel()
			before := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.RemoteBranchName{},
						RemoteSHA:  domain.SHA{},
					},
				},
			}
			after := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.NewLocalBranchName("branch-1"),
						LocalSHA:   domain.NewSHA("111111"),
						SyncStatus: domain.SyncStatusLocalOnly,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
				},
			}
			have := before.Diff(after)
			want := runstate.BranchesDiff{
				LocalAdded:   domain.LocalBranchNames{},
				LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
				LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
				RemoteAdded: []domain.RemoteBranchName{
					domain.NewRemoteBranchName("origin/branch-1"),
				},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			}
			assert.Equal(t, want, have)
		})

		t.Run("remote-only branch added", func(t *testing.T) {
			t.Parallel()
			before := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{},
			}
			after := runstate.BranchesSnapshot{
				Branches: domain.BranchInfos{
					domain.BranchInfo{
						LocalName:  domain.LocalBranchName{},
						LocalSHA:   domain.SHA{},
						SyncStatus: domain.SyncStatusRemoteOnly,
						RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
						RemoteSHA:  domain.NewSHA("111111"),
					},
				},
			}
			have := before.Diff(after)
			want := runstate.BranchesDiff{
				LocalAdded:   domain.LocalBranchNames{},
				LocalRemoved: map[domain.LocalBranchName]domain.SHA{},
				LocalChanged: map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
				RemoteAdded: []domain.RemoteBranchName{
					domain.NewRemoteBranchName("origin/branch-1"),
				},
				RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
				RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			}
			assert.Equal(t, want, have)
		})
		t.Run("remote-only branch downloaded", func(t *testing.T) {})
		t.Run("remote-only branch deleted", func(t *testing.T) {})
		t.Run("remote-only branch changed", func(t *testing.T) {})
		t.Run("local+remote branch added", func(t *testing.T) {})
		t.Run("local+remote branch removed", func(t *testing.T) {})
		t.Run("local+remote branch changed locally", func(t *testing.T) {})
		t.Run("local+remote branch updates pulled down", func(t *testing.T) {})
		t.Run("local+remote branch updates pushed up", func(t *testing.T) {})
		t.Run("local+remote branch changed remotely", func(t *testing.T) {})
		t.Run("local+remote branch changed locally and remotely to same SHA", func(t *testing.T) {})
		t.Run("local+remote branch changed locally and remotely to different SHAs", func(t *testing.T) {})
		t.Run("local+remote branch deleted locally", func(t *testing.T) {})
		t.Run("local+remote branch deleted remotely", func(t *testing.T) {})
		t.Run("local+remote branch deleted everywhere", func(t *testing.T) {})
	})
}
