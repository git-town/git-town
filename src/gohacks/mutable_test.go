package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestMutable(t *testing.T) {
	t.Parallel()
	mutable := gohacks.Mutable[gitdomain.LocalBranchNames]{}
	mutable.Get().Prepend("branch-1")
	mutable.Get().Prepend("branch-2")
	have := mutable.Get()
	want := gitdomain.NewLocalBranchNames("branch-1", "branch-2")
	must.Eq(t, &want, have)
}
