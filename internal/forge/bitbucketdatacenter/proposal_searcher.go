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

// ProposalSearcher implementation
var _ forgedomain.ProposalSearcher = bbdcConnector

func (self Connector) SearchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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
		return None[forgedomain.Proposal](), err
	}
	if len(resp.Values) == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	var needle *PullRequest
	for _, pr := range resp.Values {
		if pr.FromRef.ID == fromRefID {
			needle = &pr
			break
		}
	}
	if needle == nil {
		self.log.Success("no PR found matching source branch")
		return None[forgedomain.Proposal](), nil
	}
	proposal := parsePullRequest(*needle, self.RepositoryURL())
	self.log.Success(fmt.Sprintf("#%d", proposal.Number))
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucketDatacenter}), nil
}

func parsePullRequest(pullRequest PullRequest, repoURL string) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: false,
		Number:       pullRequest.ID,
		Source:       gitdomain.NewLocalBranchName(pullRequest.FromRef.DisplayID),
		Target:       gitdomain.NewLocalBranchName(pullRequest.ToRef.DisplayID),
		Title:        pullRequest.Title,
		Body:         NewOption(pullRequest.Description),
		URL:          fmt.Sprintf("%s/pull-requests/%v/overview", repoURL, pullRequest.ID),
	}
}
