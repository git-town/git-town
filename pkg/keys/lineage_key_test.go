package keys_test

import (
	"testing"

	"github.com/git-town/git-town/v14/pkg/keys"
	. "github.com/git-town/git-town/v14/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLineageKey(t *testing.T) {
	t.Parallel()

	t.Run("ChildName", func(t *testing.T) {
		t.Parallel()

		t.Run("valid lineage key", func(t *testing.T) {
			t.Parallel()
			key := keys.LineageKey("git-town-branch.foo.parent")
			have := key.ChildName()
			want := "foo"
			must.EqOp(t, want, have)
		})
	})

	t.Run("NewLineageKey", func(t *testing.T) {
		t.Parallel()
		tests := map[string]Option[keys.LineageKey]{
			"git-town-branch.branch.parent": Some(keys.LineageKey("git-town-branch.branch.parent")), // valid lineage key
			"git-town-branch..parent":       Some(keys.LineageKey("git-town-branch..parent")),       // empty lineage key
			"git-town.push-hook":            None[keys.LineageKey](),                                // not a lineage key
		}
		for give, want := range tests {
			key := keys.Key(give)
			have := keys.NewLineageKey(key)
			must.Eq(t, want, have)
		}
	})
}
