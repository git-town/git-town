package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumbStyle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[forgedomain.ProposalBreadcrumbStyle]
		err  bool
	}{
		{give: "", want: None[forgedomain.ProposalBreadcrumbStyle](), err: false},
		{give: "tree", want: Some(forgedomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "TREE", want: Some(forgedomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "Tree", want: Some(forgedomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "auto", want: Some(forgedomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "AUTO", want: Some(forgedomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "Auto", want: Some(forgedomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "zonk", want: None[forgedomain.ProposalBreadcrumbStyle](), err: true},
	}
	for _, tt := range tests {
		have, err := forgedomain.ParseProposalBreadcrumbStyle(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
