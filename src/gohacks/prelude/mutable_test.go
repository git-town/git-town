package prelude_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestMutable(t *testing.T) {
	t.Parallel()
	branchNames := gitdomain.LocalBranchNames{}
	mutable := Mutable[gitdomain.LocalBranchNames]{&branchNames}
	mutable.Value.Prepend("branch-1")
	mutable.Value.Prepend("branch-2")
	have := mutable.Value
	want := gitdomain.NewLocalBranchNames("branch-2", "branch-1")
	must.Eq(t, &want, have)
	must.Eq(t, want, branchNames)
}
