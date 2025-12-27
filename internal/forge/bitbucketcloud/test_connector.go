package bitbucketcloud

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

// TestConnector simulates interacting with the Bitbucket Cloud API in tests.
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
	self.log.Start(messages.APIProposalFindStart, branch, target)
	self.log.Ok()
	if self.override == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	proposal := forgedomain.BitbucketCloudProposalData{
		ProposalData: forgedomain.ProposalData{
			Active:       true,
			Body:         None[gitdomain.ProposalBody](),
			MergeWithAPI: true,
			Number:       123,
			Source:       branch,
			Target:       target,
			Title:        "title",
			URL:          self.override.String(),
		},
		CloseSourceBranch: false,
		Draft:             false,
	}
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
}
