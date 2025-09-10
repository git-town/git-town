package bitbucketcloud

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/pkg/colors"
	"github.com/ktrysmt/go-bitbucket"
)

var _ forgedomain.ProposalSourceUpdater = bbclAuthConnector

func (self AuthConnector) UpdateProposalSource(proposalData forgedomain.ProposalInterface, source gitdomain.LocalBranchName) error {
	data := proposalData.(forgedomain.BitbucketCloudProposalData)
	self.log.Start(messages.APIUpdateProposalSource, colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), colors.BoldCyan().Styled(source.String()))
	_, err := self.client.Repositories.PullRequests.Update(&bitbucket.PullRequestsOptions{
		ID:                strconv.Itoa(data.Number),
		Owner:             self.Organization,
		RepoSlug:          self.Repository,
		SourceBranch:      source.String(),
		DestinationBranch: data.Target.String(),
		Title:             data.Title,
		Description:       data.Body.GetOrZero(),
		Draft:             data.Draft,
		CloseSourceBranch: data.CloseSourceBranch,
	})
	if err != nil {
		self.log.Failed(err.Error())
		return err
	}
	self.log.Ok()
	return nil
}
