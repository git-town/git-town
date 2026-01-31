package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumb(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalBreadcrumb]
		err  bool
	}{
		{
			give: "",
			want: None[configdomain.ProposalBreadcrumb](),
			err:  false,
		},
		{
			give: "none",
			want: Some(configdomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "branches",
			want: Some(configdomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "stacks",
			want: Some(configdomain.ProposalBreadcrumbStacks),
			err:  false,
		},
		{
			give: "false",
			want: Some(configdomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "true",
			want: Some(configdomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "no",
			want: Some(configdomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "yes",
			want: Some(configdomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "0",
			want: Some(configdomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "1",
			want: Some(configdomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "zonk",
			want: None[configdomain.ProposalBreadcrumb](),
			err:  true,
		},
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalBreadcrumb(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
