package execute_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("BranchesSyncStatus", func(t *testing.T) {
		t.Run("returns the given value", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			syncStatuses := gitdomain.BranchInfos{
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
			have := fc.BranchesSyncStatus(syncStatuses, nil)
			must.Eq(t, syncStatuses, have)
			err := errors.New("test error")
			have = fc.BranchesSyncStatus(syncStatuses, err)
			must.Eq(t, syncStatuses, have)
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			fc.String("1", nil)
			fc.String("0", nil)
			must.Nil(t, fc.Err)
			fc.String("1", errors.New("first"))
			fc.String("0", errors.New("second"))
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

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string value", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			must.EqOp(t, "alpha", fc.String("alpha", nil))
			must.EqOp(t, "beta", fc.String("beta", errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			fc.String("", nil)
			must.Nil(t, fc.Err)
			fc.String("", errors.New("first"))
			fc.String("", errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string slice", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			must.Eq(t, []string{"alpha"}, fc.Strings([]string{"alpha"}, nil))
			must.Eq(t, []string{"beta"}, fc.Strings([]string{"beta"}, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := execute.FailureCollector{}
			fc.Strings([]string{}, nil)
			must.Nil(t, fc.Err)
			fc.Strings([]string{}, errors.New("first"))
			fc.Strings([]string{}, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})
}
