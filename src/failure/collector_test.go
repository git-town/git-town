package failure_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	t.Parallel()
	t.Run("Bool", func(t *testing.T) {
		t.Run("returns the given bool value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.True(t, fc.Bool(true, nil))
			assert.False(t, fc.Bool(false, nil))
			err := errors.New("test error")
			assert.True(t, fc.Bool(true, err))
			assert.False(t, fc.Bool(false, err))
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.Bool(true, nil)
			fc.Bool(false, nil)
			assert.Nil(t, fc.Err)
			fc.Bool(true, errors.New("first"))
			fc.Bool(false, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("BranchesSyncStatus", func(t *testing.T) {
		t.Run("returns the given value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			syncStatuses := git.BranchesSyncStatus{
				{
					Name:         "branch1",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
				{
					Name:         "branch2",
					InitialSHA:   "",
					SyncStatus:   git.SyncStatusLocalOnly,
					TrackingName: "",
					TrackingSHA:  "",
				},
			}
			have := fc.BranchesSyncStatus(syncStatuses, nil)
			assert.Equal(t, syncStatuses, have)
			err := errors.New("test error")
			have = fc.BranchesSyncStatus(syncStatuses, err)
			assert.Equal(t, syncStatuses, have)
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.Bool(true, nil)
			fc.Bool(false, nil)
			assert.Nil(t, fc.Err)
			fc.Bool(true, errors.New("first"))
			fc.Bool(false, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			fc := failure.Collector{}
			fc.Check(nil)
			assert.Nil(t, fc.Err)
			fc.Check(errors.New("first"))
			fc.Check(errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			fc := failure.Collector{}
			assert.False(t, fc.Check(nil))
			assert.True(t, fc.Check(errors.New("")))
			assert.True(t, fc.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			fc := failure.Collector{}
			fc.Fail("failed %s", "reason")
			assert.Error(t, fc.Err, "failed reason")
		})
	})

	t.Run("HostingService", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given HostingService value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.Equal(t, config.HostingGitHub, fc.Hosting(config.HostingGitHub, nil))
			assert.Equal(t, config.HostingGitLab, fc.Hosting(config.HostingGitLab, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.Hosting(config.HostingNone, nil)
			assert.Nil(t, fc.Err)
			fc.Hosting(config.HostingGitHub, errors.New("first"))
			fc.Hosting(config.HostingGitHub, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("PullBranchStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given PullBranchStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.Equal(t, config.PullBranchStrategyMerge, fc.PullBranchStrategy(config.PullBranchStrategyMerge, nil))
			assert.Equal(t, config.PullBranchStrategyRebase, fc.PullBranchStrategy(config.PullBranchStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, nil)
			assert.Nil(t, fc.Err)
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("first"))
			fc.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.Equal(t, "alpha", fc.String("alpha", nil))
			assert.Equal(t, "beta", fc.String("beta", errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.String("", nil)
			assert.Nil(t, fc.Err)
			fc.String("", errors.New("first"))
			fc.String("", errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string slice", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.Equal(t, []string{"alpha"}, fc.Strings([]string{"alpha"}, nil))
			assert.Equal(t, []string{"beta"}, fc.Strings([]string{"beta"}, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.Strings([]string{}, nil)
			assert.Nil(t, fc.Err)
			fc.Strings([]string{}, errors.New("first"))
			fc.Strings([]string{}, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})

	t.Run("SyncStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given SyncStrategy value", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			assert.Equal(t, config.SyncStrategyMerge, fc.SyncStrategy(config.SyncStrategyMerge, nil))
			assert.Equal(t, config.SyncStrategyRebase, fc.SyncStrategy(config.SyncStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			fc := failure.Collector{}
			fc.SyncStrategy(config.SyncStrategyMerge, nil)
			assert.Nil(t, fc.Err)
			fc.SyncStrategy(config.SyncStrategyMerge, errors.New("first"))
			fc.SyncStrategy(config.SyncStrategyMerge, errors.New("second"))
			assert.Error(t, fc.Err, "first")
		})
	})
}
