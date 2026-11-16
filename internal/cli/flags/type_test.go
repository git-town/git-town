package flags_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/cli/flags"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestBranchTypeFlag(t *testing.T) {
	t.Parallel()

	t.Run("ParseBranchTypes", func(t *testing.T) {
		t.Parallel()
		tests := map[stringss.Trimmed][]configdomain.BranchType{
			"":                      {},
			"c":                     {configdomain.BranchTypeContributionBranch},
			"c&f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c+f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c,,f":                  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c,f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c,f,o":                 {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch},
			"contribution":          {configdomain.BranchTypeContributionBranch},
			"contribution&&feature": {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution&feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution+feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution,feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"contribution|feature":  {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"c|f":                   {configdomain.BranchTypeContributionBranch, configdomain.BranchTypeFeatureBranch},
			"f":                     {configdomain.BranchTypeFeatureBranch},
			"feature":               {configdomain.BranchTypeFeatureBranch},
			"main":                  {configdomain.BranchTypeMainBranch},
			"o":                     {configdomain.BranchTypeObservedBranch},
			"observed":              {configdomain.BranchTypeObservedBranch},
			"p":                     {configdomain.BranchTypePerennialBranch},
			"pa":                    {configdomain.BranchTypeParkedBranch},
			"parked":                {configdomain.BranchTypeParkedBranch},
			"perennial":             {configdomain.BranchTypePerennialBranch},
			"pr":                    {configdomain.BranchTypePrototypeBranch},
			"prototype":             {configdomain.BranchTypePrototypeBranch},
		}
		for give, want := range tests {
			have, err := flags.ParseBranchTypes(give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})

	t.Run("SplitBranchTypeNames", func(t *testing.T) {
		t.Parallel()
		tests := map[stringss.Trimmed][]stringss.Trimmed{
			"":                   {},
			"feature":            {"feature"},
			"feature&observed":   {"feature", "observed"},
			"feature+observed":   {"feature", "observed"},
			"feature,,,observed": {"feature", "observed"},
			"feature,observed":   {"feature", "observed"},
			"feature|observed":   {"feature", "observed"},
		}
		for give, want := range tests {
			have := flags.SplitBranchTypeNames(give)
			must.Eq(t, want, have)
		}
	})
}
