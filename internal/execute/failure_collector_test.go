package execute_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v14/internal/execute"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("BranchesSyncStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given value", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			branchInfos := gitdomain.BranchInfos{
				{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch1")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				{
					LocalName:  Some(gitdomain.NewLocalBranchName("branch2")),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			have := fc.BranchInfos(branchInfos, nil)
			must.Eq(t, branchInfos, have)
			err := errors.New("test error")
			have = fc.BranchInfos(branchInfos, err)
			must.Eq(t, branchInfos, have)
		})
	})
}
