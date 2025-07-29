package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalsShowLineage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalsShowLineage]
		err  bool
	}{
		{
			give: "",
			want: None[configdomain.ProposalsShowLineage](),
			err:  false,
		},
		{
			give: "none",
			want: Some(configdomain.ProposalsShowLineageNone),
			err:  false,
		},
		{
			give: "ci",
			want: Some(configdomain.ProposalsShowLineageCI),
			err:  false,
		},
		{
			give: "CI",
			want: Some(configdomain.ProposalsShowLineageCI),
			err:  false,
		},
		{
			give: "cli",
			want: Some(configdomain.ProposalsShowLineageCLI),
			err:  false,
		},
		{
			give: "CLI",
			want: Some(configdomain.ProposalsShowLineageCLI),
			err:  false,
		},
		// {
		// 	give: "false",
		// 	want: Some(configdomain.ProposalsShowLineageNone),
		// 	err:  false,
		// },
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalsShowLineage(tt.give)
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
