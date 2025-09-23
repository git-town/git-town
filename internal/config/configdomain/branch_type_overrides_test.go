package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestBranchTypeOverrides(t *testing.T) {
	t.Parallel()

	t.Run("Concat", func(t *testing.T) {
		t.Parallel()
		data1 := configdomain.BranchTypeOverrides{
			"branch-1": configdomain.BranchTypeFeatureBranch,
			"branch-2": configdomain.BranchTypeParkedBranch,
		}
		data2 := configdomain.BranchTypeOverrides{
			"branch-1": configdomain.BranchTypeContributionBranch,
			"branch-3": configdomain.BranchTypeObservedBranch,
		}
		have := data1.Concat(data2)
		want := configdomain.BranchTypeOverrides{
			"branch-1": configdomain.BranchTypeContributionBranch,
			"branch-2": configdomain.BranchTypeParkedBranch,
			"branch-3": configdomain.BranchTypeObservedBranch,
		}
		must.Eq(t, want, have)
	})
}
