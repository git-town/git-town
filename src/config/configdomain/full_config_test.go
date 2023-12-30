package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestFullConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.True(t, config.IsFeatureBranch(gitdomain.NewLocalBranchName("feature")))
		must.False(t, config.IsFeatureBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, config.IsFeatureBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, config.IsFeatureBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("feature")))
		must.True(t, config.IsMainBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, config.IsPerennialBranch(gitdomain.NewLocalBranchName("feature")))
		must.False(t, config.IsPerennialBranch(gitdomain.NewLocalBranchName("main")))
		must.True(t, config.IsPerennialBranch(gitdomain.NewLocalBranchName("peren1")))
		must.True(t, config.IsPerennialBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
