package bitbucketcloud

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

var _ forgedomain.ProposalFinder = bbclAuthConnector

func (self AuthConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	proposalURLOverride := forgedomain.ReadProposalOverride()
	if len(proposalURLOverride) > 0 {
		return self.findProposalViaOverride(branch, target)
	}
	return self.findProposalViaAPI(branch, target)
}

func (self AuthConnector) findProposalViaAPI(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	query := fmt.Sprintf("source.branch.name = %q AND destination.branch.name = %q", branch, target)
	result1, err := self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    query,
		States:   []string{"open"},
	})
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), err
	}
	if result1 == nil {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	result2, ok := result1.(map[string]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size1, has := result2["size"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size2, ok := size1.(float64)
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	size := int(size2)
	if size == 0 {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	if size > 1 {
		self.log.Failed(fmt.Sprintf(messages.ProposalMultipleFromToFound, size, branch, target))
		return None[forgedomain.Proposal](), nil
	}
	proposal1, has := result2["values"]
	if !has {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposal2, ok := proposal1.([]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	if len(proposal2) == 0 {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposal3, ok := proposal2[0].(map[string]interface{})
	if !ok {
		self.log.Failed(messages.APIUnexpectedResultDataStructure)
		return None[forgedomain.Proposal](), nil
	}
	proposal4, err := parsePullRequest(proposal3)
	if err != nil {
		self.log.Failed(err.Error())
		return None[forgedomain.Proposal](), nil
	}
	self.log.Success(fmt.Sprintf("#%d", proposal4.Number))
	return Some(forgedomain.Proposal{Data: proposal4, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
}

func (self AuthConnector) findProposalViaOverride(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	self.log.Start(messages.APIProposalLookupStart)
	proposalURLOverride := forgedomain.ReadProposalOverride()
	self.log.Ok()
	if proposalURLOverride == forgedomain.OverrideNoProposal {
		return None[forgedomain.Proposal](), nil
	}
	proposal := forgedomain.BitbucketCloudProposalData{
		ProposalData: forgedomain.ProposalData{
			Body:         None[string](),
			MergeWithAPI: true,
			Number:       123,
			Source:       branch,
			Target:       target,
			Title:        "title",
			URL:          proposalURLOverride,
		},
		CloseSourceBranch: false,
		Draft:             false,
	}
	return Some(forgedomain.Proposal{Data: proposal, ForgeType: forgedomain.ForgeTypeBitbucket}), nil
}

func parsePullRequest(pullRequest map[string]interface{}) (result forgedomain.BitbucketCloudProposalData, err error) {
	id1, has := pullRequest["id"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	id2, ok := id1.(float64)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	number := int(id2)
	title1, has := pullRequest["title"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	title2, ok := title1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body1, has := pullRequest["description"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body2, ok := body1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination1, has := pullRequest["destination"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination2, ok := destination1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination3, has := destination2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination4, ok := destination3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination5, has := destination4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination6, ok := destination5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source1 := pullRequest["source"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source2, ok := source1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source3, has := source2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source4, ok := source3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source5, has := source4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source6, ok := source5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url1, has := pullRequest["links"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url2, ok := url1.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url3, has := url2["html"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url4, ok := url3.(map[string]interface{})
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url5, has := url4["href"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url6, ok := url5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch1, has := pullRequest["close_source_branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch2, ok := closeSourceBranch1.(bool)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft1, has := pullRequest["draft"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft2, ok := draft1.(bool)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	return forgedomain.BitbucketCloudProposalData{
		ProposalData: forgedomain.ProposalData{
			MergeWithAPI: false,
			Number:       number,
			Source:       gitdomain.NewLocalBranchName(source6),
			Target:       gitdomain.NewLocalBranchName(destination6),
			Title:        title2,
			Body:         NewOption(body2),
			URL:          url6,
		},
		CloseSourceBranch: closeSourceBranch2,
		Draft:             draft2,
	}, nil
}
