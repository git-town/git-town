package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/internal/config/configdomain"
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
}
