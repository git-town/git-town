package forgejo

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type checks
var (
	testConnector TestConnector
	_             forgedomain.Connector = testConnector
)

// TestConnector simulates interacting with the Forgejo API in tests.
type TestConnector struct {
	WebConnector
	log      print.Logger
	override forgedomain.ProposalOverride
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = testConnector // type check

func (self TestConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	self.log.Ok()
	if self.override == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         None[string](),
			MergeWithAPI: true,
			Number:       123,
			Source:       branch,
			Target:       target,
			Title:        "title",
			URL:          self.override.String(),
		},
		ForgeType: forgedomain.ForgeTypeForgejo,
	}), nil
}
