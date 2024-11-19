package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestMutable(t *testing.T) {
	t.Parallel()

	t.Run("modify the encapsulated value directly", func(t *testing.T) {
		t.Parallel()
		branchNames := gitdomain.LocalBranchNames{}
		mutable := NewMutable(&branchNames)
		mutable.Value.Prepend("branch-1")
		mutable.Value.Prepend("branch-2")
		want := gitdomain.NewLocalBranchNames("branch-2", "branch-1")
		must.Eq(t, &want, mutable.Value)
		must.Eq(t, want, branchNames)
	})

	t.Run("remains mutable when called by value", func(t *testing.T) {
		t.Parallel()
		counter := gohacks.Counter(0)
		mutable := NewMutable(&counter)
		modify(mutable)
		must.EqOp(t, 1, mutable.Immutable())
	})
}

func modify(byValue Mutable[gohacks.Counter]) {
	byValue.Value.Inc()
}
