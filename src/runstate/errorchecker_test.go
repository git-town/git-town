package runstate_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/stretchr/testify/assert"
)

func TestErrorChecker(t *testing.T) {
	t.Parallel()
	t.Run("Bool", func(t *testing.T) {
		t.Run("returns the given bool value", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.True(t, ec.Bool(true, nil))
			assert.False(t, ec.Bool(false, nil))
			err := errors.New("test error")
			assert.True(t, ec.Bool(true, err))
			assert.False(t, ec.Bool(false, err))
		})

		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.Bool(true, nil)
			ec.Bool(false, nil)
			assert.Nil(t, ec.Err)
			ec.Bool(true, errors.New("first"))
			ec.Bool(false, errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})

	t.Run("Check", func(t *testing.T) {
		t.Parallel()
		t.Run("captures the first error it receives", func(t *testing.T) {
			ec := runstate.ErrorChecker{}
			ec.Check(nil)
			assert.Nil(t, ec.Err)
			ec.Check(errors.New("first"))
			ec.Check(errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
		t.Run("indicates whether it received an error", func(t *testing.T) {
			ec := runstate.ErrorChecker{}
			assert.False(t, ec.Check(nil))
			assert.True(t, ec.Check(errors.New("")))
			assert.True(t, ec.Check(nil))
		})
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Run("registers the given error", func(t *testing.T) {
			ec := runstate.ErrorChecker{}
			ec.Fail("failed %s", "reason")
			assert.Error(t, ec.Err, "failed reason")
		})
	})

	t.Run("HostingService", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given HostingService value", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.Equal(t, config.HostingServiceGitHub, ec.HostingService(config.HostingServiceGitHub, nil))
			assert.Equal(t, config.HostingServiceGitLab, ec.HostingService(config.HostingServiceGitLab, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.HostingService(config.HostingServiceNone, nil)
			assert.Nil(t, ec.Err)
			ec.HostingService(config.HostingServiceGitHub, errors.New("first"))
			ec.HostingService(config.HostingServiceGitHub, errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})

	t.Run("PullBranchStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given PullBranchStrategy value", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.Equal(t, config.PullBranchStrategyMerge, ec.PullBranchStrategy(config.PullBranchStrategyMerge, nil))
			assert.Equal(t, config.PullBranchStrategyRebase, ec.PullBranchStrategy(config.PullBranchStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.PullBranchStrategy(config.PullBranchStrategyMerge, nil)
			assert.Nil(t, ec.Err)
			ec.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("first"))
			ec.PullBranchStrategy(config.PullBranchStrategyMerge, errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string value", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.Equal(t, "alpha", ec.String("alpha", nil))
			assert.Equal(t, "beta", ec.String("beta", errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.String("", nil)
			assert.Nil(t, ec.Err)
			ec.String("", errors.New("first"))
			ec.String("", errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given string slice", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.Equal(t, []string{"alpha"}, ec.Strings([]string{"alpha"}, nil))
			assert.Equal(t, []string{"beta"}, ec.Strings([]string{"beta"}, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.Strings([]string{}, nil)
			assert.Nil(t, ec.Err)
			ec.Strings([]string{}, errors.New("first"))
			ec.Strings([]string{}, errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})

	t.Run("SyncStrategy", func(t *testing.T) {
		t.Parallel()
		t.Run("returns the given SyncStrategy value", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			assert.Equal(t, config.SyncStrategyMerge, ec.SyncStrategy(config.SyncStrategyMerge, nil))
			assert.Equal(t, config.SyncStrategyRebase, ec.SyncStrategy(config.SyncStrategyRebase, errors.New("")))
		})
		t.Run("captures the first error it receives", func(t *testing.T) {
			t.Parallel()
			ec := runstate.ErrorChecker{}
			ec.SyncStrategy(config.SyncStrategyMerge, nil)
			assert.Nil(t, ec.Err)
			ec.SyncStrategy(config.SyncStrategyMerge, errors.New("first"))
			ec.SyncStrategy(config.SyncStrategyMerge, errors.New("second"))
			assert.Error(t, ec.Err, "first")
		})
	})
}
