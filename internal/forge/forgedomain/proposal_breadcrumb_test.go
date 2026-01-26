package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumb(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[forgedomain.ProposalBreadcrumb]
		err  bool
	}{
		{
			give: "",
			want: None[forgedomain.ProposalBreadcrumb](),
			err:  false,
		},
		{
			give: "none",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "branches",
			want: Some(forgedomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "stacks",
			want: Some(forgedomain.ProposalBreadcrumbStacks),
			err:  false,
		},
		{
			give: "false",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "true",
			want: Some(forgedomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "no",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "yes",
			want: Some(forgedomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "0",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "1",
			want: Some(forgedomain.ProposalBreadcrumbBranches),
			err:  false,
		},
		{
			give: "zonk",
			want: None[forgedomain.ProposalBreadcrumb](),
			err:  true,
		},
	}
	for _, tt := range tests {
		have, err := forgedomain.ParseProposalBreadcrumb(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
