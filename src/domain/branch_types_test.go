package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestBranchTypes(t *testing.T) {
	t.Parallel()

	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.True(t, bt.IsFeatureBranch(domain.NewLocalBranchName("feature")))
		must.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("main")))
		must.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren1")))
		must.False(t, bt.IsFeatureBranch(domain.NewLocalBranchName("peren2")))
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, bt.IsMainBranch(domain.NewLocalBranchName("feature")))
		must.True(t, bt.IsMainBranch(domain.NewLocalBranchName("main")))
		must.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren1")))
		must.False(t, bt.IsMainBranch(domain.NewLocalBranchName("peren2")))
	})

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("feature")))
		must.False(t, bt.IsPerennialBranch(domain.NewLocalBranchName("main")))
		must.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren1")))
		must.True(t, bt.IsPerennialBranch(domain.NewLocalBranchName("peren2")))
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-1", "perennial-2"),
		}
		have := branchTypes.MainAndPerennials()
		want := domain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
