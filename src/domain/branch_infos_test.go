package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfos(t *testing.T) {
	t.Run("Copy", func(t *testing.T) {
		original := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-2"),
				LocalSHA:   domain.NewSHA("222222"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		copied := original.Copy()
		original[0].LocalName = domain.NewLocalBranchName("new")
		assert.Equal(t, copied[0].LocalName, domain.NewLocalBranchName("branch-1"))
	})
}
