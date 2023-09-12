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
		bi := domain.BranchInfos{
			domain.BranchInfo{
				LocalName:  domain.NewLocalBranchName("branch-1"),
				LocalSHA:   domain.NewSHA("111111"),
				SyncStatus: domain.SyncStatusLocalOnly,
				RemoteName: domain.RemoteBranchName{},
				RemoteSHA:  domain.SHA{},
			},
		}
		cloned := bi.Clone()
		cloned[0].LocalName = domain.NewLocalBranchName("branch-2")
		cloned[0].LocalSHA = domain.NewSHA("222222")
		assert.Equal(t, bi[0].LocalName, domain.NewLocalBranchName("branch-1"))
		assert.Equal(t, bi[0].LocalSHA, domain.NewSHA("111111"))
	})
}
