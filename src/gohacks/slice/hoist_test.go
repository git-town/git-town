package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestHoist(t *testing.T) {
	t.Parallel()

	t.Run("already hoisted", func(t *testing.T) {
		t.Parallel()
		list := []string{"initial", "one", "two"}
		slice.Hoist(&list, "initial")
		want := []string{"initial", "one", "two"}
		must.Eq(t, want, list)
	})

	t.Run("contains the element to hoist", func(t *testing.T) {
		t.Parallel()
		list := []string{"alpha", "initial", "omega"}
		slice.Hoist(&list, "initial")
		want := []string{"initial", "alpha", "omega"}
		must.Eq(t, want, list)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		list := []string{}
		slice.Hoist(&list, "initial")
		want := []string{}
		must.Eq(t, want, list)
	})

	t.Run("aliased slice type", func(t *testing.T) {
		t.Parallel()
		list := domain.LocalBranchNames{domain.NewLocalBranchName("alpha"), domain.NewLocalBranchName("initial"), domain.NewLocalBranchName("omega")}
		slice.Hoist(&list, domain.NewLocalBranchName("initial"))
		want := domain.LocalBranchNames{domain.NewLocalBranchName("initial"), domain.NewLocalBranchName("alpha"), domain.NewLocalBranchName("omega")}
		must.Eq(t, want, list)
	})
}
