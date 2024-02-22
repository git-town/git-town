package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestFullConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsFeatureBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{ //nolint:exhaustruct
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
			ObservedBranches:  gitdomain.NewLocalBranchNames("observed1", "observed2"),
		}
		tests := map[string]bool{
			"feature":   true,
			"main":      false,
			"peren1":    false,
			"peren2":    false,
			"peren3":    true,
			"observed1": false,
			"observed2": false,
			"observed3": true,
		}
		for give, want := range tests {
			have := config.IsFeatureBranch(gitdomain.NewLocalBranchName(give))
			must.Eq(t, want, have)
		}
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{ //nolint:exhaustruct
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
		config := configdomain.FullConfig{ //nolint:exhaustruct
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
			PerennialRegex:    "release-.*",
		}
		tests := map[string]bool{
			"main":      false,
			"peren1":    true,
			"peren2":    true,
			"peren3":    false,
			"feature":   false,
			"release-1": true,
			"release-2": true,
			"other":     false,
		}
		for give, want := range tests {
			have := config.IsPerennialBranch(gitdomain.NewLocalBranchName(give))
			must.Eq(t, want, have)
		}
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		config := configdomain.FullConfig{ //nolint:exhaustruct
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
