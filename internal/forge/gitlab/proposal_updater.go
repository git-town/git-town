package gitlab

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/pkg/colors"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (self AuthConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription string) error {
	data := proposalData.Data()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)))
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), data.Number, &gitlab.UpdateMergeRequestOptions{
		Description: Ptr(updatedDescription),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}

func (self AuthConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	data := proposalData.Data()
	self.log.Start(messages.ForgeGitLabUpdateMRViaAPI, data.Number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), data.Number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.Ptr(target.String()),
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}
