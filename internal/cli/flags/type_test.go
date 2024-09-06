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
			"":  {},
			"f": {configdomain.BranchTypeFeatureBranch},
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
