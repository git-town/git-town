package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseBreadcrumbExcludeBranches(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalBreadcrumbExcludeBranches]
		err  bool
	}{
		{
			give: "",
			want: Some(configdomain.NewProposalBreadcrumbExcludeBranches()),
			err:  false,
		},
		{
			give: "prototype,contribution",
			want: Some(configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: " prototype, contribution ",
			want: Some(configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: "prototype,prototype",
			want: Some(configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch)),
			err:  false,
		},
		{
			give: "p,c",
			want: Some(configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: "zonk",
			want: None[configdomain.ProposalBreadcrumbExcludeBranches](),
			err:  true,
		},
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalBreadcrumbExcludeBranches(stringss.Trim(tt.give), "test")
		must.EqOp(t, tt.err, err != nil)
		must.True(t, have.Equal(tt.want))
	}
}

func TestParseBreadcrumbExcludeBranchesList(t *testing.T) {
	t.Parallel()
	have, err := configdomain.ParseProposalBreadcrumbExcludeBranchesList([]string{"prototype", " contribution ", ""}, "test")
	must.NoError(t, err)
	want := Some(configdomain.NewProposalBreadcrumbExcludeBranches(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch))
	must.True(t, have.Equal(want))
}
