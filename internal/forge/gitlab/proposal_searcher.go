package gitlab

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (self AuthConnector) SearchProposalFn() Option[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	if self.APIToken.IsNone() {
		return None[func(gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
	}
	return Some(self.searchProposal)
}

func (self AuthConnector) searchProposal(branch gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIParentBranchLookupStart, branch.String())
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	case 1:
		proposal := parseMergeRequest(mergeRequests[0])
		self.log.Success(proposal.Target.String())
		return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeGitLab}), nil
	default:
		return None[forgedomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFromFound, len(mergeRequests), branch)
	}
}
