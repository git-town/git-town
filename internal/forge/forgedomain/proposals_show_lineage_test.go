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
		want Option[forgedomain.ProposalsShowLineage]
		err  bool
	}{
		{
			give: "",
			want: None[forgedomain.ProposalsShowLineage](),
			err:  false,
		},
		{
			give: "none",
			want: Some(forgedomain.ProposalsShowLineageNone),
			err:  false,
		},
		{
			give: "ci",
			want: Some(forgedomain.ProposalsShowLineageCI),
			err:  false,
		},
		{
			give: "CI",
			want: Some(forgedomain.ProposalsShowLineageCI),
			err:  false,
		},
		{
			give: "Ci",
			want: Some(forgedomain.ProposalsShowLineageCI),
			err:  false,
		},
		{
			give: "cli",
			want: Some(forgedomain.ProposalsShowLineageCLI),
			err:  false,
		},
		{
			give: "false",
			want: Some(forgedomain.ProposalsShowLineageNone),
			err:  false,
		},
		{
			give: "true",
			want: Some(forgedomain.ProposalsShowLineageCLI),
			err:  false,
		},
		{
			give: "no",
			want: Some(forgedomain.ProposalsShowLineageNone),
			err:  false,
		},
		{
			give: "yes",
			want: Some(forgedomain.ProposalsShowLineageCLI),
			err:  false,
		},
		{
			give: "0",
			want: Some(forgedomain.ProposalsShowLineageNone),
			err:  false,
		},
		{
			give: "1",
			want: Some(forgedomain.ProposalsShowLineageCLI),
			err:  false,
		},
		{
			give: "zonk",
			want: None[forgedomain.ProposalsShowLineage](),
			err:  true,
		},
	}
	for _, tt := range tests {
		have, err := forgedomain.ParseProposalsShowLineage(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
