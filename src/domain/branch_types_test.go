package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/shoenig/test"
)

func TestBranchTypes(t *testing.T) {
	t.Parallel()

	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		test.True(t, bt.IsFeatureBranch(domain.NewLocalBranchName("feature")))
		test.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("main")))
		test.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren1")))
		test.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren2")))
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		test.False(t, bt.IsMainBranch(domain.NewLocalBranchName("feature")))
		test.True(t, bt.IsMainBranch(domain.NewLocalBranchName("main")))
		test.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren1")))
		test.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren2")))
	})

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		test.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("feature")))
		test.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("main")))
		test.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren1")))
		test.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren2")))
	})
}
