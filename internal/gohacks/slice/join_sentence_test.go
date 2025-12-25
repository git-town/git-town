package slice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/shoenig/test/must"
)

func TestJoinSentenceQuotes(t *testing.T) {
	t.Parallel()

	t.Run("0 elements", func(t *testing.T) {
		t.Parallel()
		give := []configdomain.BranchType{}
		have := slice.JoinSentenceQuotes(give)
		want := ""
		must.EqOp(t, want, have)
	})

	t.Run("1 element", func(t *testing.T) {
		t.Parallel()
		give := []configdomain.BranchType{configdomain.BranchTypeMainBranch}
		have := slice.JoinSentenceQuotes(give)
		want := `"main"`
		must.EqOp(t, want, have)
	})

	t.Run("2 elements", func(t *testing.T) {
		t.Parallel()
		give := []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypeFeatureBranch}
		have := slice.JoinSentenceQuotes(give)
		want := `"main" and "feature"`
		must.EqOp(t, want, have)
	})

	t.Run("3 elements", func(t *testing.T) {
		t.Parallel()
		give := []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch}
		have := slice.JoinSentenceQuotes(give)
		want := `"main", "feature", and "contribution"`
		must.EqOp(t, want, have)
	})

	t.Run("4 elements", func(t *testing.T) {
		t.Parallel()
		give := []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch}
		have := slice.JoinSentenceQuotes(give)
		want := `"main", "feature", "contribution", and "observed"`
		must.EqOp(t, want, have)
	})
}
