package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchTypes(t *testing.T) {
	t.Parallel()

	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.True(t, bt.IsFeatureBranch(gitdomain.NewLocalBranchName("feature")))
		must.False(t, bt.IsFeatureBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, bt.IsFeatureBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, bt.IsFeatureBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, bt.IsMainBranch(gitdomain.NewLocalBranchName("feature")))
		must.True(t, bt.IsMainBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, bt.IsMainBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, bt.IsMainBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		bt := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, bt.IsPerennialBranch(gitdomain.NewLocalBranchName("feature")))
		must.False(t, bt.IsPerennialBranch(gitdomain.NewLocalBranchName("main")))
		must.True(t, bt.IsPerennialBranch(gitdomain.NewLocalBranchName("peren1")))
		must.True(t, bt.IsPerennialBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
		}
		have := branchTypes.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
