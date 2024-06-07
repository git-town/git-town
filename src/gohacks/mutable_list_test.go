package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestMutableList(t *testing.T) {
	t.Parallel()

	t.Run("empty mutable list", func(t *testing.T) {
		t.Parallel()
		list := gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]{}
		have := list.List()
		want := gitdomain.LocalBranchNames{}
		must.Eq(t, want, have)
	})
	t.Run("populate mutable list directly", func(t *testing.T) {
		t.Parallel()
		list := gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]{}
		list.Append(gitdomain.NewLocalBranchName("branch-1"))
		list.Append(gitdomain.NewLocalBranchName("branch-2"))
		have := list.List()
		want := gitdomain.NewLocalBranchNames("branch-1", "branch-2")
		must.Eq(t, want, have)
	})
	t.Run("populate mutable list given as argument", func(t *testing.T) {
		t.Parallel()
		list := gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]{}
		branch1 := gitdomain.NewLocalBranchName("branch-1")
		branch2 := gitdomain.NewLocalBranchName("branch-2")
		appendBranch(list, branch1)
		appendBranch(list, branch2)
		have := list.List()
		want := gitdomain.LocalBranchNames{branch1, branch2}
		must.Eq(t, want, have)
	})
}

type branchList = gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]

func appendBranch(list branchList, branch gitdomain.LocalBranchName) {
	list.Append(branch)
}
