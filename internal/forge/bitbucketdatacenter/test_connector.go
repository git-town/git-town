package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	testConnector TestConnector
	_             forgedomain.Connector = testConnector
)

// TestConnector simulates interacting with the Bitbucket DataCenter API in tests.
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
	data := forgedomain.ProposalData{
		Body:         None[string](),
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          self.override.String(),
	}
	return Some(forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}
