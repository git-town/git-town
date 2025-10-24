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
			have, err := configdomain.ParseDisplayTypes("all", "unit test")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierAll,
				BranchTypes: []configdomain.BranchType{},
			}
			must.True(t, have.EqualSome(want))
		})
		t.Run("no", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayTypes("no", "unit test")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{},
			}
			must.True(t, have.EqualSome(want))
		})

		t.Run("exclude branch types", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayTypes("no feature prototype", "unit test")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			must.True(t, have.EqualSome(want))
		})

		t.Run("only branch types", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayTypes("observed contribution", "unit test")
			must.NoError(t, err)
			want := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierOnly,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeObservedBranch, configdomain.BranchTypeContributionBranch},
			}
			must.True(t, have.EqualSome(want))
		})

		t.Run("empty string", func(t *testing.T) {
			t.Parallel()
			have, err := configdomain.ParseDisplayTypes("", "unit test")
			must.Nil(t, err)
			must.True(t, have.IsNone())
		})

		t.Run("invalid string", func(t *testing.T) {
			t.Parallel()
			_, err := configdomain.ParseDisplayTypes("zonk", "unit test")
			must.NotNil(t, err)
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
					t.Parallel()
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
					t.Parallel()
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
					t.Parallel()
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
			want := "all branch types"
			must.EqOp(t, want, have)
		})

		t.Run("no", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{},
			}
			have := give.String()
			want := "no branch types"
			must.EqOp(t, want, have)
		})

		t.Run("exclude specific branches", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierNo,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeMainBranch},
			}
			have := give.String()
			want := `all branch types except "feature", "prototype", and "main"`
			must.EqOp(t, want, have)
		})

		t.Run("only specific branches", func(t *testing.T) {
			t.Parallel()
			give := configdomain.DisplayTypes{
				Quantifier:  configdomain.QuantifierOnly,
				BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch},
			}
			have := give.String()
			want := `only the branch types "feature" and "prototype"`
			must.EqOp(t, want, have)
		})
	})
}
