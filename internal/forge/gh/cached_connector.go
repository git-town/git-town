package gh

import (
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// type checks
var (
	cachedConnector CachedConnector
	_               forgedomain.Connector = &cachedConnector
)

// CachedConnector talks to the GitHub API through the "gh" executable with caching.
type CachedConnector struct {
	Cache     forgedomain.APICache
	Connector Connector
}

func (self *CachedConnector) BrowseRepository(runner subshelldomain.Runner) error {
	return self.Connector.BrowseRepository(runner)
}

func (self *CachedConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	return self.Connector.CreateProposal(data)
}

func (self *CachedConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return self.Connector.DefaultProposalMessage(data)
}

func (self *CachedConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal, has := self.Cache.Lookup(source, target); has {
		return cachedProposal, nil
	}
	loadedProposal, err := self.Connector.FindProposal(source, target)
	if err == nil {
		self.Cache.RegisterLookupResult(source, target, loadedProposal)
	}
	return loadedProposal, err
}

func (self *CachedConnector) ProposalReference(data forgedomain.ProposalData) string {
	return self.Connector.ProposalReference(data)
}

func (self *CachedConnector) SearchProposals(source gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedSearchResult, has := self.Cache.LookupSearch(source).Get(); has {
		return cachedSearchResult, nil
	}
	loadedSearchResult, err := self.Connector.SearchProposals(source)
	if err == nil {
		self.Cache.RegisterSearchResult(source, loadedSearchResult)
	}
	return loadedSearchResult, err
}

func (self *CachedConnector) SquashMergeProposal(number forgedomain.ProposalNumber, message Option[gitdomain.CommitMessage]) error {
	self.Cache.Clear(number)
	return self.Connector.SquashMergeProposal(number, message)
}

func (self *CachedConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedBody gitdomain.ProposalBody) error {
	self.Cache.Clear(proposalData.Data().Number)
	return self.Connector.UpdateProposalBody(proposalData, updatedBody)
}

func (self *CachedConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.Cache.Clear(proposalData.Data().Number)
	return self.Connector.UpdateProposalTarget(proposalData, target)
}

func (self *CachedConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	return self.Connector.VerifyCredentials()
}

// ============================================================================
// browse the repo
// ============================================================================

// ============================================================================
// create proposals
// ============================================================================

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &cachedConnector // type-check

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &cachedConnector // type-check

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = &cachedConnector // type-check

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &cachedConnector // type-check

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &cachedConnector // type-check

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = &cachedConnector // type check
