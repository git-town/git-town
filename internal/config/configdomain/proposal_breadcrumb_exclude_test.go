package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseBreadcrumbExclude(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalBreadcrumbExclude]
		err  bool
	}{
		{
			give: "",
			want: Some(configdomain.NewProposalBreadcrumbExclude()),
			err:  false,
		},
		{
			give: "prototype contribution",
			want: Some(configdomain.NewProposalBreadcrumbExclude(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: " prototype contribution ",
			want: Some(configdomain.NewProposalBreadcrumbExclude(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: "p c",
			want: Some(configdomain.NewProposalBreadcrumbExclude(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeContributionBranch)),
			err:  false,
		},
		{
			give: "prototype, contribution",
			want: None[configdomain.ProposalBreadcrumbExclude](),
			err:  true,
		},
		{
			give: "zonk",
			want: None[configdomain.ProposalBreadcrumbExclude](),
			err:  true,
		},
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalBreadcrumbExclude(stringss.Trim(tt.give), "test")
		must.EqOp(t, tt.err, err != nil)
		must.True(t, have.Equal(tt.want))
	}
}

func TestParseBreadcrumbExcludeList(t *testing.T) {
	t.Parallel()
	have, err := configdomain.ParseProposalBreadcrumbExcludeList([]string{"prototype", " contribution ", ""}, "test")
	must.NoError(t, err)
	want := Some(configdomain.NewProposalBreadcrumbExclude(configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeContributionBranch))
	must.True(t, have.Equal(want))
}
