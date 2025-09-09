package bitbucketdatacenter

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to enforce conformance to the ProposalFinder interface
var _ forgedomain.ProposalFinder = bbdcAPIConnector

func (self APIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	proposalURLOverride := forgedomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		return self.findProposalViaOverride(branch, target)
	}
	return self.findProposalViaAPI(branch, target)
}

func (self APIConnector) apiBaseURL() string {
	return fmt.Sprintf(
		"https://%s/rest/api/latest/projects/%s/repos/%s/pull-requests",
		self.Hostname,
		self.Organization,
		self.Repository,
	)
}

func (self APIConnector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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
	proposal := parsePullRequest(*needle, self.RepositoryURL())
	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}

func (self APIConnector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := forgedomain.ReadProposalOverride()
	self.log.Ok()
	if proposalURLOverride == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	data := forgedomain.ProposalData{
		Body:         None[string](),
		MergeWithAPI: true,
		Number:       123,
		Source:       branch,
		Target:       target,
		Title:        "title",
		URL:          proposalURLOverride,
	}
	return Some(forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}
