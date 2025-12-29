package github

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
	Log      print.Logger
	Override forgedomain.ProposalOverride
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = testConnector // type check

func (self TestConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.Log.Start(messages.APIProposalFindStart, branch, target)
	self.Log.Ok()
	if self.Override == forgedomain.OverrideNoProposal {
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
			URL:          self.Override.String(),
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}), nil
}
