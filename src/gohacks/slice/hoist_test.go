package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestHoist(t *testing.T) {
	t.Parallel()
	t.Run("already hoisted", func(t *testing.T) {
		t.Parallel()
		give := []string{"initial", "one", "two"}
		have := slice.Hoist(give, "initial")
		want := []string{"initial", "one", "two"}
		must.Eq(t, want, have)
	})
	t.Run("contains the element to hoist", func(t *testing.T) {
		t.Parallel()
		give := []string{"alpha", "initial", "omega"}
		have := slice.Hoist(give, "initial")
		want := []string{"initial", "alpha", "omega"}
		must.Eq(t, want, have)
	})
	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		have := slice.Hoist(give, "initial")
		want := []string{}
		must.Eq(t, want, have)
	})
	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		give := domain.LocalBranchNames{domain.NewLocalBranchName("alpha"), domain.NewLocalBranchName("initial"), domain.NewLocalBranchName("omega")}
		have := slice.Hoist(give, domain.NewLocalBranchName("initial"))
		want := domain.LocalBranchNames{domain.NewLocalBranchName("initial"), domain.NewLocalBranchName("alpha"), domain.NewLocalBranchName("omega")}
		must.Eq(t, want, have)
	})
}
