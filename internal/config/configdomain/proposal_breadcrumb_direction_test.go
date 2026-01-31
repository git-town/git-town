package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseProposalBreadcrumbDirection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		give string
		want Option[configdomain.ProposalBreadcrumbDirection]
		err  bool
	}{
		{give: "", want: None[configdomain.ProposalBreadcrumbDirection](), err: false},
		{give: "down", want: Some(configdomain.ProposalBreadcrumbDirectionDown), err: false},
		{give: "DOWN", want: Some(configdomain.ProposalBreadcrumbDirectionDown), err: false},
		{give: "Down", want: Some(configdomain.ProposalBreadcrumbDirectionDown), err: false},
		{give: "up", want: Some(configdomain.ProposalBreadcrumbDirectionUp), err: false},
		{give: "UP", want: Some(configdomain.ProposalBreadcrumbDirectionUp), err: false},
		{give: "Up", want: Some(configdomain.ProposalBreadcrumbDirectionUp), err: false},
		{give: "zonk", want: None[configdomain.ProposalBreadcrumbDirection](), err: true},
	}
	for _, tt := range tests {
		have, err := configdomain.ParseProposalBreadcrumbDirection(tt.give, "test")
		must.EqOp(t, err != nil, tt.err)
		must.True(t, have.Equal(tt.want))
	}
}
