package config_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
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
	have, err := config.NewBranchTypeOverridesInSnapshot(snapshot, false, nil)
	must.NoError(t, err)
	want := configdomain.BranchTypeOverrides{
		"branch-1": configdomain.BranchTypeFeatureBranch,
		"branch-2": configdomain.BranchTypeObservedBranch,
	}
	must.Eq(t, want, have)
}

func TestNewPartialConfigFromSnapshotProposalBreadcrumbExclude(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		snapshot := configdomain.SingleSnapshot{
			configdomain.KeyProposalBreadcrumbExclude: "",
		}
		have, err := config.NewPartialConfigFromSnapshot(snapshot, false, false, nil)
		must.NoError(t, err)
		want := configdomain.NewProposalBreadcrumbExclude()
		must.True(t, have.ProposalBreadcrumbExclude.EqualSome(want))
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		snapshot := configdomain.SingleSnapshot{
			configdomain.KeyProposalBreadcrumbExclude: "zonk",
		}
		_, err := config.NewPartialConfigFromSnapshot(snapshot, false, false, nil)
		must.Error(t, err)
	})

	t.Run("set", func(t *testing.T) {
		t.Parallel()
		snapshot := configdomain.SingleSnapshot{
			configdomain.KeyProposalBreadcrumbExclude: "prototype contribution",
		}
		have, err := config.NewPartialConfigFromSnapshot(snapshot, false, false, nil)
		must.NoError(t, err)
		want := configdomain.NewProposalBreadcrumbExclude(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch)
		must.True(t, have.ProposalBreadcrumbExclude.EqualSome(want))
	})
}
