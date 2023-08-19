package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchDurations(t *testing.T) {
	t.Parallel()
	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.True(t, bt.IsFeatureBranch(domain.NewLocalBranchName("feature")))
		assert.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("main")))
		assert.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren1")))
		assert.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren2")))
	})
	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.False(t, bt.IsMainBranch(domain.NewLocalBranchName("feature")))
		assert.True(t, bt.IsMainBranch(domain.NewLocalBranchName("main")))
		assert.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren1")))
		assert.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren2")))
	})
	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("feature")))
		assert.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("main")))
		assert.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren1")))
		assert.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren2")))
	})
}
