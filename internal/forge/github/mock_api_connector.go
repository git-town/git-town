package github

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

// type checks
var (
	mockAPIConnector MockAPIConnector
	_                forgedomain.Connector = &mockAPIConnector
)

// MockAPIConnector provides access to the Bitbucket Cloud API while caching proposal information.
type MockAPIConnector struct {
	WebConnector
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &mockAPIConnector // type check

func (self *MockAPIConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal, has := self.cache.Lookup(source, target); has {
		return cachedProposal, nil
	}
	loadedProposal, err := self.api.FindProposal(source, target)
	if err == nil {
		self.cache.RegisterLookupResult(source, target, loadedProposal)
	}
	return loadedProposal, err
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &mockAPIConnector // type check

func (self *MockAPIConnector) SearchProposals(source gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedSearchResult, has := self.cache.LookupSearch(source); has {
		return cachedSearchResult, nil
	}
	loadedSearchResult, err := self.api.SearchProposals(source)
	if err == nil {
		self.cache.RegisterSearchResult(source, loadedSearchResult)
	}
	return loadedSearchResult, err
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = &mockAPIConnector // type check

func (self *MockAPIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	self.cache.Clear()
	return self.api.SquashMergeProposal(number, message)
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &mockAPIConnector // type check

func (self *MockAPIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody gitdomain.ProposalBody) error {
	self.cache.Clear()
	return self.api.UpdateProposalBody(proposalData, newBody)
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &mockAPIConnector // type check

func (self *MockAPIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.cache.Clear()
	return self.api.UpdateProposalTarget(proposalData, target)
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = &mockAPIConnector // type check

func (self *MockAPIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	return self.api.VerifyCredentials()
}
