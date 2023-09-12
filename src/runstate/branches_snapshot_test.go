package runstate_test

import (
	"testing"
)

func TestBranchesSnapshot(t *testing.T) {
	t.Parallel()
	t.Run("Diff", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch added", func(t *testing.T) {
			t.Parallel()
			// before := runstate.BranchesSnapshot{
			// 	Branches: domain.BranchInfos{},
			// }
			// after := runstate.BranchesSnapshot{
			// 	Branches: domain.BranchInfos{
			// 		domain.BranchInfo{
			// 			LocalName:  domain.NewLocalBranchName("branch-1"),
			// 			LocalSHA:   domain.NewSHA("111111"),
			// 			SyncStatus: domain.SyncStatusLocalOnly,
			// 			RemoteName: domain.RemoteBranchName{},
			// 			RemoteSHA:  domain.SHA{},
			// 		},
			// 	},
			// }
			// have := before.Diff(after)
			// want := runstate.BranchesDiff{
			// 	LocalAdded:    domain.NewLocalBranchNames("branch-1"),
			// 	LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
			// 	LocalChanged:  map[domain.LocalBranchName]runstate.Change[domain.SHA]{},
			// 	RemoteAdded:   []domain.RemoteBranchName{},
			// 	RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
			// 	RemoteChanged: map[domain.RemoteBranchName]runstate.Change[domain.SHA]{},
			// }
			// assert.Equal(t, want, have)
		})
	})
}
