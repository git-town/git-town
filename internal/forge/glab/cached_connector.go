package glab

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	cachedConnector CachedConnector
	_               forgedomain.Connector = cachedConnector
)

// CachedConnector provides standardized connectivity for the given repository
// via the glab CLI with caching.
type CachedConnector struct {
	Connector Connector
	Cache     forgedomain.ProposalCache
}

// ============================================================================
// browse the repo
// ============================================================================

func (self CachedConnector) BrowseRepository(runner subshelldomain.Runner) error {
	return self.Connector.BrowseRepository(runner)
}

// ============================================================================
// create proposals
// ============================================================================

func (self CachedConnector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	return self.Connector.CreateProposal(data)
}

func (self CachedConnector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return self.Connector.DefaultProposalMessage(data)
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = cachedConnector // type check

func (self CachedConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal := self.Cache.BySourceTarget(branch, target); cachedProposal.IsSome() {
		return cachedProposal, nil
	}
	result, err := self.Connector.FindProposal(branch, target)
	self.Cache.SetOption(result)
	return result, err
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = cachedConnector // type check

func (self CachedConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedProposals := self.Cache.BySource(branch); len(cachedProposals) > 0 {
		return cachedProposals, nil
	}
	result, err := self.Connector.SearchProposals(branch)
	self.Cache.SetMany(result)
	return result, err
}

// ============================================================================
// squash-merge proposals
// ============================================================================

var _ forgedomain.ProposalMerger = cachedConnector // type check

func (self CachedConnector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	self.Cache.Delete(number)
	return self.Connector.SquashMergeProposal(number, message)
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = cachedConnector // type check

func (self CachedConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, updatedDescription string) error {
	self.Cache.Delete(proposalData.Data().Number)
	return self.Connector.UpdateProposalBody(proposalData, updatedDescription)
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = cachedConnector // type check

func (self CachedConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.Cache.Delete(proposalData.Data().Number)
	return self.Connector.UpdateProposalTarget(proposalData, target)
}

// ============================================================================
// verify credentials
// ============================================================================

var _ forgedomain.CredentialVerifier = cachedConnector

func (self CachedConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	return self.Connector.VerifyCredentials()
}
