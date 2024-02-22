package configdomain_test

import (
	"fmt"
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
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			ObservedBranches:  gitdomain.NewLocalBranchNames("observed-1", "observed-2"),
			ParkedBranches:    gitdomain.NewLocalBranchNames("parked-1", "parked-2"),
		}
		tests := map[string]bool{
			"feature":     true,
			"main":        false,
			"perennial-1": false,
			"perennial-2": false,
			"perennial-3": true,
			"observed-1":  false,
			"observed-2":  false,
			"observed-3":  true,
			"parked-1":    false,
			"parked-2":    false,
			"parked-3":    true,
		}
		for give, want := range tests {
			have := config.IsFeatureBranch(gitdomain.NewLocalBranchName(give))
			fmt.Println(give)
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
