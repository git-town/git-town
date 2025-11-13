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
	_                  forgedomain.Connector = cachedAPIConnector
)

// APIConnector provides access to the Bitbucket DataCenter API.
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

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = cachedAPIConnector // type check

func (self CachedAPIConnector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
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

var _ forgedomain.ProposalSearcher = cachedAPIConnector // type check

func (self CachedAPIConnector) SearchProposals(branch gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedProposals := self.cache.BySource(branch); len(cachedProposals) > 0 {
		return cachedProposals, nil
	}
	result, err := self.api.SearchProposals(branch)
	self.cache.SetMany(result)
	return result, err
}
