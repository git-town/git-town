package execute_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("BranchesSyncStatus", func(t *testing.T) {
		t.Run("returns the given value", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			branchInfos := gitdomain.BranchInfos{
				{
					LocalName:  gitdomain.NewLocalBranchName("branch1"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
				{
					LocalName:  gitdomain.NewLocalBranchName("branch2"),
					LocalSHA:   gitdomain.EmptySHA(),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: gitdomain.EmptyRemoteBranchName(),
					RemoteSHA:  gitdomain.EmptySHA(),
				},
			}
			have := fc.BranchInfos(branchInfos, nil)
			must.Eq(t, branchInfos, have)
			err := errors.New("test error")
			have = fc.BranchInfos(branchInfos, err)
			must.Eq(t, branchInfos, have)
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			fc := execute.FailureCollector{}
			fc.Check(nil)
			must.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			fc := execute.FailureCollector{}
			must.False(t, fc.Check(nil))
			must.True(t, fc.Check(errors.New("")))
			must.True(t, fc.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			fc := execute.FailureCollector{}
			fc.Fail("failed %s", "reason")
			must.ErrorContains(t, fc.Err, "failed reason")
		})
	})
}
