package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseShareNewBranches(t *testing.T) {
	t.Parallel()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		have, err := configdomain.ParseShareNewBranches("", "test source")
		must.NoError(t, err)
		must.True(t, have.IsNone())
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()
		_, err := configdomain.ParseShareNewBranches("zonk", "test source")
		must.Error(t, err)
		must.StrContains(t, err.Error(), "invalid value for \"test source\": \"zonk\"")
	})

	t.Run("valid values", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.ShareNewBranches]{
			"no":      Some(configdomain.ShareNewBranchesNone),
			"false":   Some(configdomain.ShareNewBranchesNone),
			"0":       Some(configdomain.ShareNewBranchesNone),
			"push":    Some(configdomain.ShareNewBranchesPush),
			"propose": Some(configdomain.ShareNewBranchesPropose),
		}
		for give, want := range tests {
			t.Run(give, func(t *testing.T) {
				t.Parallel()
				have, err := configdomain.ParseShareNewBranches(give, "test source")
				must.NoError(t, err)
				must.Eq(t, want, have)
			})
		}
	})
}
