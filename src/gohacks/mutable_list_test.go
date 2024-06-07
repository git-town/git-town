package gohacks_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/shoenig/test/must"
)

func TestMutableList(t *testing.T) {
	t.Parallel()

	t.Run(".List", func(t *testing.T) {
		t.Run("empty mutable list", func(t *testing.T) {
			t.Parallel()
			list := gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]{}
			have := list.List()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})
		t.Run("populated mutable list", func(t *testing.T) {
			t.Parallel()
			list := gohacks.MutableList[gitdomain.LocalBranchName, gitdomain.LocalBranchNames]{}
			list.Append(gitdomain.NewLocalBranchName("branch-1"))
			list.Append(gitdomain.NewLocalBranchName("branch-2"))
			have := list.List()
			want := gitdomain.NewLocalBranchNames("branch-1", "branch-2")
			must.Eq(t, want, have)
		})
	})
}
