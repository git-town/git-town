package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumbStyle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalBreadcrumbStyle]
		err  bool
	}{
		{give: "", want: None[configdomain.ProposalBreadcrumbStyle](), err: false},
		{give: "tree", want: Some(configdomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "TREE", want: Some(configdomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "Tree", want: Some(configdomain.ProposalBreadcrumbStyleTree), err: false},
		{give: "auto", want: Some(configdomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "AUTO", want: Some(configdomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "Auto", want: Some(configdomain.ProposalBreadcrumbStyleAuto), err: false},
		{give: "zonk", want: None[configdomain.ProposalBreadcrumbStyle](), err: true},
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalBreadcrumbStyle(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
