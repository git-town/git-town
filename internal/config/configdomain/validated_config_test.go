package configdomain_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsMainOrPerennialBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.UnvalidatedConfig{
			ContributionBranches: gitdomain.NewLocalBranchNames("contribution"),
			MainBranch:           Some(gitdomain.NewLocalBranchName("main")),
			PerennialBranches:    gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			ObservedBranches:     gitdomain.NewLocalBranchNames("observed"),
			ParkedBranches:       gitdomain.NewLocalBranchNames("parked"),
		}
		tests := map[string]bool{
			"feature":     false,
			"main":        true,
			"perennial-1": true,
			"perennial-2": true,
			"perennial-3": false,
			"observed":    false,
			"parked":      false,
		}
		for give, want := range tests {
			have := config.IsMainOrPerennialBranch(gitdomain.NewLocalBranchName(give))
			fmt.Println(give)
			must.Eq(t, want, have)
		}
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.UnvalidatedConfig{
			MainBranch:        Some(gitdomain.NewLocalBranchName("main")),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
		}
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("feature")))
		must.True(t, config.IsMainBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.UnvalidatedConfig{
			MainBranch:        Some(gitdomain.NewLocalBranchName("main")),
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
			PerennialRegex:    configdomain.ParsePerennialRegex("release-.*"),
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
		config := configdomain.UnvalidatedConfig{
			MainBranch:        Some(gitdomain.NewLocalBranchName("main")),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
