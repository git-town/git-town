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
	_                  forgedomain.Connector = &cachedAPIConnector
)

// CachedAPIConnector provides access to the Bitbucket Cloud API while caching proposal information.
type CachedAPIConnector struct {
	api   APIConnector
	cache forgedomain.ProposalCache
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

func (self *CachedAPIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal := self.cache.BySourceTarget(branch, target); cachedProposal.IsSome() {
		return cachedProposal, nil
	}
	result, err := self.api.FindProposal(branch, target)
	self.cache.SetOption(result)
	return result, err
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &cachedAPIConnector // type check

func (self *CachedAPIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedProposals := self.cache.BySource(branch); len(cachedProposals) > 0 {
		return cachedProposals, nil
	}
	result, err := self.api.SearchProposals(branch)
	self.cache.SetMany(result)
	return result, err
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = &cachedAPIConnector // type check

func (self *CachedAPIConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	self.cache.Delete(number)
	return self.api.SquashMergeProposal(number, message)
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &cachedAPIConnector // type check

func (self *CachedAPIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody string) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalBody(proposalData, newBody)
}

// ============================================================================
// udpate proposal source
// ============================================================================

var _ forgedomain.ProposalSourceUpdater = &cachedAPIConnector // type check

func (self *CachedAPIConnector) UpdateProposalSource(proposalData forgedomain.ProposalInterface, source gitdomain.LocalBranchName) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalSource(proposalData, source)
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &cachedAPIConnector // type check

func (self *CachedAPIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.cache.Delete(proposalData.Data().Number)
	return self.api.UpdateProposalTarget(proposalData, target)
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = &cachedAPIConnector // type check

func (self *CachedAPIConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	return self.api.VerifyCredentials()
}
