package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestAncestry(t *testing.T) {
	t.Parallel()

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesWithParentAndSyncStatus{
			git.BranchWithParentAndSyncStatus{
				Name:       "up-to-date",
				SyncStatus: git.SyncStatusUpToDate,
			},
			git.BranchWithParentAndSyncStatus{
				Name:       "ahead",
				SyncStatus: git.SyncStatusAhead,
			},
			git.BranchWithParentAndSyncStatus{
				Name:       "behind",
				SyncStatus: git.SyncStatusBehind,
			},
			git.BranchWithParentAndSyncStatus{
				Name:       "local-only",
				SyncStatus: git.SyncStatusLocalOnly,
			},
			git.BranchWithParentAndSyncStatus{
				Name:       "remote-only",
				SyncStatus: git.SyncStatusRemoteOnly,
			},
			git.BranchWithParentAndSyncStatus{
				Name:       "deleted-at-remote",
				SyncStatus: git.SyncStatusDeletedAtRemote,
			},
		}
		want := []string{"up-to-date", "ahead", "behind", "local-only"}
		have := bs.LocalBranches().BranchNames()
		assert.Equal(t, want, have)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		bs := git.BranchesWithParentAndSyncStatus{
			git.BranchWithParentAndSyncStatus{
				Name: "one",
			},
			git.BranchWithParentAndSyncStatus{
				Name: "two",
			},
		}
		assert.Equal(t, "one", bs.Lookup("one").Name)
		assert.Equal(t, "two", bs.Lookup("two").Name)
		assert.Nil(t, bs.Lookup("zonk"))
	})

}
