package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalsShowLineage(t *testing.T) {
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
			give: "ci",
			want: Some(forgedomain.ProposalsBreadcrumbCI),
			err:  false,
		},
		{
			give: "CI",
			want: Some(forgedomain.ProposalsBreadcrumbCI),
			err:  false,
		},
		{
			give: "Ci",
			want: Some(forgedomain.ProposalsBreadcrumbCI),
			err:  false,
		},
		{
			give: "cli",
			want: Some(forgedomain.ProposalBreadcrumbCLI),
			err:  false,
		},
		{
			give: "false",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "true",
			want: Some(forgedomain.ProposalBreadcrumbCLI),
			err:  false,
		},
		{
			give: "no",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "yes",
			want: Some(forgedomain.ProposalBreadcrumbCLI),
			err:  false,
		},
		{
			give: "0",
			want: Some(forgedomain.ProposalBreadcrumbNone),
			err:  false,
		},
		{
			give: "1",
			want: Some(forgedomain.ProposalBreadcrumbCLI),
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
