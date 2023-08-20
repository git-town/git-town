package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfos(t *testing.T) {
	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Run("there is a remote branch matching the given local branch", func(t *testing.T) {
			bis := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.LocalBranchName{},
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.NewRemoteBranchName("origin/branch-1"),
					RemoteSHA:  domain.SHA{},
				},
			}
			give := domain.NewLocalBranchName("branch-1")
			assert.True(t, bis.HasMatchingRemoteBranchFor(give))
		})
		t.Run("there is a local branch matching the given local branch", func(t *testing.T) {
			bis := domain.BranchInfos{
				domain.BranchInfo{
					Name:       domain.NewLocalBranchName("branch-1"),
					InitialSHA: domain.SHA{},
					SyncStatus: domain.SyncStatusUpToDate,
					RemoteName: domain.RemoteBranchName{},
					RemoteSHA:  domain.SHA{},
				},
			}
			give := domain.NewLocalBranchName("branch-1")
			assert.False(t, bis.HasMatchingRemoteBranchFor(give))
		})
	})
}
