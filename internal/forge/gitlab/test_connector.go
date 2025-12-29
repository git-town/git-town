package gitlab

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

// TestConnector simulates interacting with the GitHub API in tests.
type TestConnector struct {
	WebConnector
	log      print.Logger
	override forgedomain.ProposalOverride
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = testConnector

func (self TestConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalFindStart, branch, target)
	self.log.Ok()
	if self.override == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Active:       true,
			Body:         None[gitdomain.ProposalBody](),
			MergeWithAPI: true,
			Number:       123,
			Source:       branch,
			Target:       target,
			Title:        "title",
			URL:          self.override.String(),
		},
		ForgeType: forgedomain.ForgeTypeGitLab,
	}), nil
}
