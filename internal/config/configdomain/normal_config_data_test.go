package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestNormalConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsPerennialBranch", func(t *testing.T) {
		t.Parallel()
		perennialRegexOpt, err := configdomain.ParsePerennialRegex("release-.*")
		must.NoError(t, err)
		config := configdomain.NormalConfigData{
			PerennialBranches: gitdomain.NewLocalBranchNames("peren1", "peren2"),
			PerennialRegex:    perennialRegexOpt,
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

	t.Run("IsPrototypeBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("listed as prototype branch", func(t *testing.T) {
			t.Parallel()
			config := configdomain.NormalConfigData{
				PrototypeBranches: gitdomain.NewLocalBranchNames("proto1"),
				DefaultBranchType: configdomain.BranchTypeFeatureBranch,
			}
		})
	})
}
