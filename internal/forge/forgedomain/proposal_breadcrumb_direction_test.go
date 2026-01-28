package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumbDirection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[forgedomain.ProposalBreadcrumbDirection]
		err  bool
	}{
		{give: "", want: None[forgedomain.ProposalBreadcrumbDirection](), err: false},
		{give: "top-down", want: Some(forgedomain.ProposalBreadcrumbDirectionTopDown), err: false},
		{give: "bottom-up", want: Some(forgedomain.ProposalBreadcrumbDirectionBottomUp), err: false},
		{give: "zonk", want: None[forgedomain.ProposalBreadcrumbDirection](), err: true},
	}
	for _, tt := range tests {
		have, err := forgedomain.ParseProposalBreadcrumbDirection(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
