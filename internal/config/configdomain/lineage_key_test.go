package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLineageKey(t *testing.T) {
	t.Parallel()

	t.Run("ChildBranch", func(t *testing.T) {
		t.Parallel()
		key := configdomain.NewLineageKey("git-town-branch.my-branch.parent")
		have := key.ChildBranch()
		want := gitdomain.LocalBranchName("my-branch")
		must.EqOp(t, want, have)
	})

	t.Run("ParseLineageKey", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.LineageKey]{
			"git-town-branch.branch.parent": Some(configdomain.NewLineageKey("git-town-branch.branch.parent")), // valid lineage key
			"git-town-branch..parent":       Some(configdomain.NewLineageKey("git-town-branch..parent")),       // empty lineage key
			"git-town.push-hook":            None[configdomain.LineageKey](),                                   // not a lineage key
		}
		for give, want := range tests {
			key := configdomain.Key(give)
			have := configdomain.ParseLineageKey(key)
			must.Eq(t, want, have)
		}
	})
}
