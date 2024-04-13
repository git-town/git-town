package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestHoist(t *testing.T) {
	t.Parallel()

	t.Run("already hoisted", func(t *testing.T) {
		t.Parallel()
		list := []string{"initial", "one", "two"}
		have := slice.Hoist(list, "initial")
		want := []string{"initial", "one", "two"}
		must.Eq(t, want, have)
	})

	t.Run("contains the element to hoist", func(t *testing.T) {
		t.Parallel()
		list := []string{"alpha", "initial", "omega"}
		have := slice.Hoist(list, "initial")
		want := []string{"initial", "alpha", "omega"}
		must.Eq(t, want, have)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		list := []string{}
		have := slice.Hoist(list, "initial")
		want := []string{}
		must.Eq(t, want, have)
	})

	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		list := gitdomain.LocalBranchNames{gitdomain.NewLocalBranchName("alpha"), gitdomain.NewLocalBranchName("initial"), gitdomain.NewLocalBranchName("omega")}
		have := slice.Hoist(list, gitdomain.NewLocalBranchName("initial"))
		want := gitdomain.LocalBranchNames{gitdomain.NewLocalBranchName("initial"), gitdomain.NewLocalBranchName("alpha"), gitdomain.NewLocalBranchName("omega")}
		must.Eq(t, want, have)
	})
}
