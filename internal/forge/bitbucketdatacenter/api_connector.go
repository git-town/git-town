package bitbucketdatacenter

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	apiConnector APIConnector
	_            forgedomain.Connector = apiConnector
)

// APIConnector provides access to the Bitbucket DataCenter API.
type APIConnector struct {
	WebConnector
	cache    forgedomain.ProposalCache
	log      print.Logger
	token    string
	username string
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = apiConnector // type check

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if proposal := self.cache.BySourceTarget(branch, target); proposal.IsSome() {
		return proposal, nil
	}
	self.log.Start(messages.APIProposalLookupStart)
	ctx := context.TODO()
	fromRefID := fmt.Sprintf("refs/heads/%v", branch)
	toRefID := fmt.Sprintf("refs/heads/%v", target)
	var resp PullRequestResponse
	err := requests.URL(self.apiBaseURL()).
		BasicAuth(self.username, self.token).
		Param("at", toRefID).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	var needle *PullRequest
	for _, pr := range resp.Values {
		if pr.FromRef.ID == fromRefID && pr.ToRef.ID == toRefID {
			needle = &pr
			break
		}
	}
	if needle == nil {
		self.log.Success("no PR found matching source and target branch")
		return None[forgedomain.Proposal](), nil
	}
	proposalData := parsePullRequest(*needle, self.RepositoryURL())
	self.log.Success(fmt.Sprintf("#%d", proposalData.Number))
	proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}
	self.cache.Set(proposal)
	return Some(proposal), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = apiConnector // type check

func (self APIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if proposal := self.cache.BySourceTarget(branch, target); proposal.IsSome() {
		return proposal, nil
	}
	self.log.Start(messages.APIProposalLookupStart)
	ctx := context.TODO()
	fromRefID := fmt.Sprintf("refs/heads/%v", branch)
	var resp PullRequestResponse
	err := requests.URL(self.apiBaseURL()).
		BasicAuth(self.username, self.token).
		ToJSON(&resp).
		Fetch(ctx)
	if err != nil {
		self.log.Failed(err.Error())
		return []forgedomain.Proposal{}, err
	}
	result := []forgedomain.Proposal{}
	for _, pr := range resp.Values {
		if pr.FromRef.ID == fromRefID {
			proposalData := parsePullRequest(pr, self.RepositoryURL())
			self.log.Success(fmt.Sprintf("#%d ", proposalData.Number))
			proposal := forgedomain.Proposal{Data: proposalData, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}
			result = append(result, proposal)
		}
	}
	if len(result) == 0 {
		self.log.Success("none")
	}
	return result, nil
}

func (self APIConnector) apiBaseURL() string {
	return fmt.Sprintf(
		"https://%s/rest/api/latest/projects/%s/repos/%s/pull-requests",
		self.Hostname,
		self.Organization,
		self.Repository,
	)
}
