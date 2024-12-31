package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
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

	t.Run("NewBranchTypeOverridesFromSnapshot", func(t *testing.T) {
		t.Parallel()
		snapshot := configdomain.SingleSnapshot{
			"git-town-branch.branch-1.branchtype": "feature",
			"git-town-branch.branch-2.branchtype": "observed",
			"git-town-branch.branch-3.parent":     "main",
			"git-town.prototype-branches":         "foo",
		}
		removeFunc := func(configdomain.Key) error { return nil }
		have, err := configdomain.NewBranchTypeOverridesInSnapshot(snapshot, removeFunc)
		must.NoError(t, err)
		want := configdomain.BranchTypeOverrides{
			"branch-1": configdomain.BranchTypeFeatureBranch,
			"branch-2": configdomain.BranchTypeObservedBranch,
		}
		must.Eq(t, want, have)
	})
}
