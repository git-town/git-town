package gohacks_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("Bool", func(t *testing.T) {
		t.Run("returns the given bool value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			must.True(t, fc.Bool(true, nil))
			must.False(t, fc.Bool(false, nil))
			err := errors.New("test error")
			must.True(t, fc.Bool(true, err))
			must.False(t, fc.Bool(false, err))
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			fc.Bool(true, nil)
			fc.Bool(false, nil)
			must.Nil(t, fc.Err)
			fc.Bool(true, errors.New("first"))
			fc.Bool(false, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("BranchesSyncStatus", func(t *testing.T) {
		t.Run("returns the given value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			syncStatuses := domain.BranchInfos{
				{
					LocalName:  domain.NewLocalBranchName("branch1"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
				},
				{
					LocalName:  domain.NewLocalBranchName("branch2"),
					LocalSHA:   domain.EmptySHA(),
					SyncStatus: domain.SyncStatusLocalOnly,
					RemoteName: domain.EmptyRemoteBranchName(),
					RemoteSHA:  domain.EmptySHA(),
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
			fc := gohacks.FailureCollector{}
			fc.Bool(true, nil)
			fc.Bool(false, nil)
			must.Nil(t, fc.Err)
			fc.Bool(true, errors.New("first"))
			fc.Bool(false, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			fc := gohacks.FailureCollector{}
			fc.Check(nil)
			must.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			fc := gohacks.FailureCollector{}
			must.False(t, fc.Check(nil))
			must.True(t, fc.Check(errors.New("")))
			must.True(t, fc.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			fc := gohacks.FailureCollector{}
			fc.Fail("failed %s", "reason")
			must.ErrorContains(t, fc.Err, "failed reason")
		})
	})

	t.Run("HostingService", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given HostingService value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			must.EqOp(t, config.HostingGitHub, fc.Hosting(config.HostingGitHub, nil))
			must.EqOp(t, config.HostingGitLab, fc.Hosting(config.HostingGitLab, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			fc.Hosting(config.HostingNone, nil)
			must.Nil(t, fc.Err)
			fc.Hosting(config.HostingGitHub, errors.New("first"))
			fc.Hosting(config.HostingGitHub, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("PullBranchStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given PullBranchStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			must.EqOp(t, config.PullBranchStrategyMerge, fc.PullBranchStrategy(config.PullBranchStrategyMerge, nil))
			must.EqOp(t, config.PullBranchStrategyRebase, fc.PullBranchStrategy(config.PullBranchStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, nil)
			must.Nil(t, fc.Err)
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("first"))
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			must.EqOp(t, "alpha", fc.String("alpha", nil))
			must.EqOp(t, "beta", fc.String("beta", errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
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
			fc := gohacks.FailureCollector{}
			must.Eq(t, []string{"alpha"}, fc.Strings([]string{"alpha"}, nil))
			must.Eq(t, []string{"beta"}, fc.Strings([]string{"beta"}, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			fc.Strings([]string{}, nil)
			must.Nil(t, fc.Err)
			fc.Strings([]string{}, errors.New("first"))
			fc.Strings([]string{}, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("SyncStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given SyncStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			must.EqOp(t, config.SyncStrategyMerge, fc.SyncStrategy(config.SyncStrategyMerge, nil))
			must.EqOp(t, config.SyncStrategyRebase, fc.SyncStrategy(config.SyncStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := gohacks.FailureCollector{}
			fc.SyncStrategy(config.SyncStrategyMerge, nil)
			must.Nil(t, fc.Err)
			fc.SyncStrategy(config.SyncStrategyMerge, errors.New("first"))
			fc.SyncStrategy(config.SyncStrategyMerge, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})
}
