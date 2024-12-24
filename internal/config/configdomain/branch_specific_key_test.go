package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestBranchSpecificKey(t *testing.T) {

	t.Run("ChildName", func(t *testing.T) {
		t.Parallel()
		key := configdomain.NewLineageKey("git-town-branch.foo.parent")
		have := key.BranchName()
		want := "foo"
		must.EqOp(t, want, have)
	})
}
