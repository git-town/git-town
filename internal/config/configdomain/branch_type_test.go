package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchType(t *testing.T) {
	t.Parallel()

	t.Run("ParseBranchType", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.BranchType]{
			"":             None[configdomain.BranchType](),
			"(none)":       None[configdomain.BranchType](),
			"contribution": Some(configdomain.BranchTypeContributionBranch),
			"feature":      Some(configdomain.BranchTypeFeatureBranch),
			"main":         Some(configdomain.BranchTypeMainBranch),
			"observed":     Some(configdomain.BranchTypeObservedBranch),
			"parked":       Some(configdomain.BranchTypeParkedBranch),
			"perennial":    Some(configdomain.BranchTypePerennialBranch),
			"prototype":    Some(configdomain.BranchTypePrototypeBranch),
			"f":            Some(configdomain.BranchTypeFeatureBranch),
			"fe":           Some(configdomain.BranchTypeFeatureBranch),
			"fea":          Some(configdomain.BranchTypeFeatureBranch),
			"pa":           Some(configdomain.BranchTypeParkedBranch),
			"pe":           Some(configdomain.BranchTypePerennialBranch),
			"pr":           Some(configdomain.BranchTypePrototypeBranch),
		}
		for give, want := range tests {
			have, err := configdomain.ParseBranchType(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})
}
