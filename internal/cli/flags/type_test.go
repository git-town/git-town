package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestBranchTypeFlag(t *testing.T) {
	t.Parallel()

	t.Run("ParseBranchTypes", func(t *testing.T) {
		t.Parallel()
		tests := map[string][]configdomain.BranchType{
			"":                      {},
			"contribution":          {configdomain.BranchTypeContributionBranch},
			"feature":               {configdomain.BranchTypeFeatureBranch},
			"observed":              {configdomain.BranchTypeObservedBranch},
			"perennial":             {configdomain.BranchTypePerennialBranch},
			"parked":                {configdomain.BranchTypeParkedBranch},
			"prototype":             {configdomain.BranchTypePrototypeBranch},
			"main":                  {configdomain.BranchTypeMainBranch},
			"contribution,feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution+feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution&feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution|feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution&&feature": {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c":                     {configdomain.BranchTypeContributionBranch},
			"f":                     {configdomain.BranchTypeFeatureBranch},
			"o":                     {configdomain.BranchTypeObservedBranch},
			"p":                     {configdomain.BranchTypePerennialBranch},
			"pa":                    {configdomain.BranchTypeParkedBranch},
			"pr":                    {configdomain.BranchTypePrototypeBranch},
			"c,f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c+f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c&f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c|f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c,,f":                  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c,f,o":                 {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch},
		}
		for give, want := range tests {
			have, err := flags.ParseBranchTypes(give)
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("SplitBranchTypeNames", func(t *testing.T) {
		t.Parallel()
		tests := map[string][]string{
			"":                   {},
			"feature":            {"feature"},
			"feature,observed":   {"feature", "observed"},
			"feature|observed":   {"feature", "observed"},
			"feature&observed":   {"feature", "observed"},
			"feature+observed":   {"feature", "observed"},
			"feature,,,observed": {"feature", "observed"},
		}
		for give, want := range tests {
			have := flags.SplitBranchTypeNames(give)
			must.Eq(t, want, have)
		}
	})
}
