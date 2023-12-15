package configdomain_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	t.Run("Bool", func(t *testing.T) {
		t.Run("returns the given bool value", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			must.True(t, fc.Bool(true, nil))
			must.False(t, fc.Bool(false, nil))
			err := errors.New("test error")
			must.True(t, fc.Bool(true, err))
			must.False(t, fc.Bool(false, err))
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
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
			fc := configdomain.FailureCollector{}
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
			fc := configdomain.FailureCollector{}
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
			fc := configdomain.FailureCollector{}
			fc.Check(nil)
			must.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			fc := configdomain.FailureCollector{}
			must.False(t, fc.Check(nil))
			must.True(t, fc.Check(errors.New("")))
			must.True(t, fc.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			fc := configdomain.FailureCollector{}
			fc.Fail("failed %s", "reason")
			must.ErrorContains(t, fc.Err, "failed reason")
		})
	})

	t.Run("HostingService", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given HostingService value", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			must.EqOp(t, configdomain.HostingGitHub, fc.Hosting(configdomain.HostingGitHub, nil))
			must.EqOp(t, configdomain.HostingGitLab, fc.Hosting(configdomain.HostingGitLab, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			fc.Hosting(configdomain.HostingNone, nil)
			must.Nil(t, fc.Err)
			fc.Hosting(configdomain.HostingGitHub, errors.New("first"))
			fc.Hosting(configdomain.HostingGitHub, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("SyncPerennialStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given SyncPerennialStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			must.EqOp(t, configdomain.SyncPerennialStrategyMerge, fc.SyncPerennialStrategy(configdomain.SyncPerennialStrategyMerge, nil))
			must.EqOp(t, configdomain.SyncPerennialStrategyRebase, fc.SyncPerennialStrategy(configdomain.SyncPerennialStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			fc.SyncPerennialStrategy(configdomain.SyncPerennialStrategyMerge, nil)
			must.Nil(t, fc.Err)
			fc.SyncPerennialStrategy(configdomain.SyncPerennialStrategyMerge, errors.New("first"))
			fc.SyncPerennialStrategy(configdomain.SyncPerennialStrategyMerge, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string value", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			must.EqOp(t, "alpha", fc.String("alpha", nil))
			must.EqOp(t, "beta", fc.String("beta", errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
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
			fc := configdomain.FailureCollector{}
			must.Eq(t, []string{"alpha"}, fc.Strings([]string{"alpha"}, nil))
			must.Eq(t, []string{"beta"}, fc.Strings([]string{"beta"}, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			fc.Strings([]string{}, nil)
			must.Nil(t, fc.Err)
			fc.Strings([]string{}, errors.New("first"))
			fc.Strings([]string{}, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})

	t.Run("SyncFeatureStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given SyncFeatureStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			must.EqOp(t, configdomain.SyncFeatureStrategyMerge, fc.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, nil))
			must.EqOp(t, configdomain.SyncFeatureStrategyRebase, fc.SyncFeatureStrategy(configdomain.SyncFeatureStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := configdomain.FailureCollector{}
			fc.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, nil)
			must.Nil(t, fc.Err)
			fc.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, errors.New("first"))
			fc.SyncFeatureStrategy(configdomain.SyncFeatureStrategyMerge, errors.New("second"))
			must.ErrorContains(t, fc.Err, "first")
		})
	})
}
