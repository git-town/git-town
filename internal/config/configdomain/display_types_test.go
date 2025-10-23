package configdomain_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestDisplayTypes(t *testing.T) {
	t.Parallel()

	t.Run("ParseDisplayType", func(t *testing.T) {
		t.Parallel()

		t.Run("all", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayType("all")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierAll,
				BranchTypes: []configdomain.BranchType{},
			}
			must.Eq(t, want, have)
		})
		t.Run("no", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayType("no")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{},
			}
			must.Eq(t, want, have)
		})

		t.Run("exclude branch types", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayType("no feature prototype")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			must.Eq(t, want, have)
		})

		t.Run("only branch types", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayType("observed contribution")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierOnly,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("ShouldDisplayType", func(t *testing.T) {
		t.Parallel()
		t.Run("all", func(t *testing.T) {
			t.Parallel()
			displayTypes := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierAll,
				BranchTypes: []configdomain.BranchType{},
			}
			for _, branchType := range configdomain.AllBranchTypes() {
				must.True(t, displayTypes.ShouldDisplayType(branchType))
			}
		})
		t.Run("no", func(t *testing.T) {
			t.Parallel()
			displayTypes := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{},
			}
			for _, branchType := range configdomain.AllBranchTypes() {
				t.Run(branchType.String(), func(t *testing.T) {
					must.False(t, displayTypes.ShouldDisplayType(branchType))
				})
			}
		})
		t.Run("exclude specific types", func(t *testing.T) {
			t.Parallel()
			displayTypes := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			tests := map[configdomain.BranchType]bool{
				configdomain.BranchTypeFeatureBranch:      false,
				configdomain.BranchTypePrototypeBranch:    false,
				configdomain.BranchTypePerennialBranch:    true,
				configdomain.BranchTypeMainBranch:         true,
				configdomain.BranchTypeObservedBranch:     true,
				configdomain.BranchTypeContributionBranch: true,
			}
			for give, want := range tests {
				t.Run(fmt.Sprintf("%s ==> %t", give, want), func(t *testing.T) {
					have := displayTypes.ShouldDisplayType(give)
					must.EqOp(t, want, have)
				})
			}
		})
		t.Run("only specific types", func(t *testing.T) {
			t.Parallel()
			displayTypes := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierOnly,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch},
			}
			tests := map[configdomain.BranchType]bool{
				configdomain.BranchTypeFeatureBranch:      false,
				configdomain.BranchTypePrototypeBranch:    false,
				configdomain.BranchTypePerennialBranch:    false,
				configdomain.BranchTypeMainBranch:         false,
				configdomain.BranchTypeObservedBranch:     true,
				configdomain.BranchTypeContributionBranch: true,
			}
			for give, want := range tests {
				t.Run(fmt.Sprintf("%s ==> %t", give, want), func(t *testing.T) {
					have := displayTypes.ShouldDisplayType(give)
					must.EqOp(t, want, have)
				})
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		t.Run("all", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierAll,
				BranchTypes: []configdomain.BranchType{},
			}
			have := give.String()
			want := "all"
			must.EqOp(t, want, have)
		})

		t.Run("no", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{},
			}
			have := give.String()
			want := "no"
			must.EqOp(t, want, have)
		})

		t.Run("exclude specific branches", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			have := give.String()
			want := "no feature prototype"
			must.EqOp(t, want, have)
		})

		t.Run("only specific branches", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierOnly,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			have := give.String()
			want := "feature prototype"
			must.EqOp(t, want, have)
		})
	})
}
