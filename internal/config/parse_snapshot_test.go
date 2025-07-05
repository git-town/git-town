package config_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/shoenig/test/must"
)

func TestNewBranchTypeOverridesFromSnapshot(t *testing.T) {
	t.Parallel()
	snapshot := configdomain.SingleSnapshot{
		"git-town-branch.branch-1.branchtype": "feature",
		"git-town-branch.branch-2.branchtype": "observed",
		"git-town-branch.branch-3.parent":     "main",
		"git-town.prototype-branches":         "foo",
	}
	gitIO := gitconfig.IO{Shell: nil}
	have, err := config.NewBranchTypeOverridesInSnapshot(snapshot, &gitIO)
	must.NoError(t, err)
	want := configdomain.BranchTypeOverrides{
		"branch-1": configdomain.BranchTypeFeatureBranch,
		"branch-2": configdomain.BranchTypeObservedBranch,
	}
	must.Eq(t, want, have)
}
