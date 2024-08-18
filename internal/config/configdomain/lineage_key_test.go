package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configdomain"
	. "github.com/git-town/git-town/v15/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLineageKey(t *testing.T) {
	t.Parallel()

	t.Run("ChildName", func(t *testing.T) {
		t.Parallel()

		t.Run("valid lineage key", func(t *testing.T) {
			t.Parallel()
			key := configdomain.LineageKey("git-town-branch.foo.parent")
			have := key.ChildName()
			want := "foo"
			must.EqOp(t, want, have)
		})
	})

	t.Run("NewLineageKey", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[configdomain.LineageKey]{
			"git-town-branch.branch.parent": Some(configdomain.LineageKey("git-town-branch.branch.parent")), // valid lineage key
			"git-town-branch..parent":       Some(configdomain.LineageKey("git-town-branch..parent")),       // empty lineage key
			"git-town.push-hook":            None[configdomain.LineageKey](),                                // not a lineage key
		}
		for give, want := range tests {
			key := configdomain.Key(give)
			have := configdomain.NewLineageKey(key)
			must.Eq(t, want, have)
		}
	})
}
