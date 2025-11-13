package bitbucketdatacenter

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	cachedAPIConnector CachedAPIConnector
	_                  forgedomain.Connector = &cachedAPIConnector
)

// APIConnector provides access to the Bitbucket DataCenter API.
type CachedAPIConnector struct {
	api   APIConnector
	cache forgedomain.APICache
}

func (self *CachedAPIConnector) BrowseRepository(runner subshelldomain.Runner) error {
	return self.api.BrowseRepository(runner)
}

func (self *CachedAPIConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	return self.api.CreateProposal(data)
}

func (self *CachedAPIConnector) DefaultProposalMessage(proposalData forgedomain.ProposalData) string {
	return self.api.DefaultProposalMessage(proposalData)
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &cachedAPIConnector // type check

func (self *CachedAPIConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal, has := self.cache.Lookup(source, target); has {
		return cachedProposal, nil
	}
	loadedProposal, err := self.api.FindProposal(source, target)
	if err != nil {
		self.cache.RegisterLookupResult(source, target, loadedProposal)
	}
	return loadedProposal, err
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &cachedAPIConnector // type check

func (self *CachedAPIConnector) SearchProposals(source gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedSearchResult, has := self.cache.LookupSearch(source); has {
		return cachedSearchResult, nil
	}
	loadedSearchResult, err := self.api.SearchProposals(source)
	if err != nil {
		self.cache.RegisterSearchResult(source, loadedSearchResult)
	}
	return loadedSearchResult, err
}
