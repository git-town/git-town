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
		bd := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.True(t, bd.IsFeatureBranch(domain.NewLocalBranchName("feature")))
		assert.False(t, bd.IsFeatureBranch(domain.NewLocalBranchName("main")))
		assert.False(t, bd.IsFeatureBranch(domain.NewLocalBranchName("peren1")))
		assert.False(t, bd.IsFeatureBranch(domain.NewLocalBranchName("peren2")))
	})
	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bd := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.False(t, bd.IsMainBranch(domain.NewLocalBranchName("feature")))
		assert.True(t, bd.IsMainBranch(domain.NewLocalBranchName("main")))
		assert.False(t, bd.IsMainBranch(domain.NewLocalBranchName("peren1")))
		assert.False(t, bd.IsMainBranch(domain.NewLocalBranchName("peren2")))
	})
	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bd := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		assert.False(t, bd.IsPerennialBranch(domain.NewLocalBranchName("feature")))
		assert.False(t, bd.IsPerennialBranch(domain.NewLocalBranchName("main")))
		assert.True(t, bd.IsPerennialBranch(domain.NewLocalBranchName("peren1")))
		assert.True(t, bd.IsPerennialBranch(domain.NewLocalBranchName("peren2")))
	})
}
