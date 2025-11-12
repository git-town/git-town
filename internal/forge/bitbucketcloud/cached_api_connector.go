package bitbucketcloud

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	cachedAPIConnector CachedAPIConnector
	_                  forgedomain.Connector = cachedAPIConnector
)

// APIConnector provides access to the Bitbucket Cloud API.
type CachedAPIConnector struct {
	api   APIConnector
	cache forgedomain.ProposalCache
}

func (self CachedAPIConnector) BrowseRepository(runner subshelldomain.Runner) error {
	return self.api.BrowseRepository(runner)
}

func (self CachedAPIConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	return self.api.CreateProposal(data)
}

func (self CachedAPIConnector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	return self.api.DefaultProposalMessage(proposalData)
}

func (self CachedAPIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal := self.cache.BySourceTarget(branch, target); cachedProposal.IsSome() {
		return cachedProposal, nil
	}
	result, err := self.api.FindProposal(branch, target)
	self.cache.SetOption(result)
	return result, err
}

func (self CachedAPIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedProposals := self.cache.BySource(branch); len(cachedProposals) > 0 {
		return cachedProposals, nil
	}
	result, err := self.api.SearchProposals(branch)
	self.cache.SetMany(result)
	return result, err
}

func (self CachedAPIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	self.cache.Delete(number)
	return self.api.SquashMergeProposal(number, message)
}

func (self CachedAPIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody string) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalBody(proposalData, newBody)
}

func (self CachedAPIConnector) UpdateProposalSource(proposalData forgedomain.ProposalInterface, source gitdomain.LocalBranchName) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalSource(proposalData, source)
}

func (self CachedAPIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalTarget(proposalData, target)
}

func (self CachedAPIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	return self.api.VerifyCredentials()
}
