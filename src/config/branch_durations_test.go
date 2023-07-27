package config_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/stretchr/testify/assert"
)

func TestBranchDurations(t *testing.T) {
	t.Parallel()
	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		bd := config.BranchDurations{
			MainBranch:        "main",
			PerennialBranches: []string{"peren1", "peren2"},
		}
		assert.True(t, bd.IsFeatureBranch("feature"))
		assert.False(t, bd.IsFeatureBranch("main"))
		assert.False(t, bd.IsFeatureBranch("peren1"))
		assert.False(t, bd.IsFeatureBranch("peren2"))
	})
	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bd := config.BranchDurations{
			MainBranch:        "main",
			PerennialBranches: []string{"peren1", "peren2"},
		}
		assert.False(t, bd.IsMainBranch("feature"))
		assert.True(t, bd.IsMainBranch("main"))
		assert.False(t, bd.IsMainBranch("peren1"))
		assert.False(t, bd.IsMainBranch("peren2"))
	})
	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bd := config.BranchDurations{
			MainBranch:        "main",
			PerennialBranches: []string{"peren1", "peren2"},
		}
		assert.False(t, bd.IsPerennialBranch("feature"))
		assert.False(t, bd.IsPerennialBranch("main"))
		assert.True(t, bd.IsPerennialBranch("peren1"))
		assert.True(t, bd.IsPerennialBranch("peren2"))
	})
}
