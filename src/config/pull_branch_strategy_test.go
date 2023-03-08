package config_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/stretchr/testify/assert"
)

func TestNewPullBranchStrategy(t *testing.T) {
	t.Parallel()
	t.Run("valid content", func(t *testing.T) {
		t.Parallel()
		tests := map[string]config.PullBranchStrategy{
			"merge":  config.PullBranchStrategyMerge,
			"rebase": config.PullBranchStrategyRebase,
		}
		for give, want := range tests {
			have, err := config.NewPullBranchStrategy(give)
			assert.Nil(t, err)
			assert.Equal(t, want, have)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		t.Parallel()
		for _, give := range []string{"merge", "Merge", "MERGE"} {
			have, err := config.NewPullBranchStrategy(give)
			assert.Nil(t, err)
			assert.Equal(t, config.PullBranchStrategyMerge, have)
		}
	})

	t.Run("defaults to rebase", func(t *testing.T) {
		t.Parallel()
		have, err := config.NewPullBranchStrategy("")
		assert.Nil(t, err)
		assert.Equal(t, config.PullBranchStrategyRebase, have)
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()
		_, err := config.NewPullBranchStrategy("zonk")
		assert.Error(t, err)
	})
}
