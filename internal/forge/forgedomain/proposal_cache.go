package forgedomain

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCache caches known proposals and knowledge about non-existing proposals.
type ProposalCache struct {
	// the cached items
	results []Result
}

type Result interface {
}

// the result of a find operation
type lookupResult struct {
	source   gitdomain.LocalBranchName
	target   gitdomain.LocalBranchName
	proposal Option[Proposal]
}

// the result of a search operation
type searchResult struct {
	source    gitdomain.LocalBranchName
	Proposals []Proposal
}

// Lookup provides what this cache knows about the proposal for the given source and target branch.
// If it has a cached proposal, returns (Some, true).
// If it knowns that there is no proposal, returns (None, true).
// If it has never heard about the given source and target branch, returns (None, false).
func (self *ProposalCache) Lookup(source, target gitdomain.LocalBranchName) (proposal Option[Proposal], knows bool) {
	for _, result := range self.results {
		switch result := result.(type) {
		case lookupResult:
			if result.source == source && result.target == target {
				return result.proposal, true
			}
		case searchResult:
			if result.source == source {
				for _, proposal := range result.Proposals {
					if proposal.Data.Data().Target == target {
						return Some(proposal), true
					}
				}
				// here we know that there was a proposal search for the source branch,
				// and the target branch was not in the result --> we know there is no proposal
				return None[Proposal](), true
			}
		}
		panic(fmt.Sprintf("unknown result type: %T", result))
	}
	// we didn't run across any API results for the source branch
	return None[Proposal](), false
}

// LookupSearch provides the cached search result for the given source branch.
func (self *ProposalCache) LookupSearch(source gitdomain.LocalBranchName) (proposals []Proposal, knows bool) {
	for _, result := range self.results {
		if searchResult, ok := result.(searchResult); ok {
			if searchResult.source == source {
				return searchResult.Proposals, true
			}
		}
	}
	return []Proposal{}, false
}

// Clear removes all cached results.
func (self *ProposalCache) Clear() {
	self.results = []Result{}
}

// SaveLookupResult registers the given result of a lookup operation.
func (self *ProposalCache) RegisterLookupResult(source, target gitdomain.LocalBranchName, proposal Option[Proposal]) {
	self.removeLookupResult(source, target)
	self.results = append(self.results, lookupResult{
		source:   source,
		target:   target,
		proposal: proposal,
	})
}

// RegisterSearchResult registers the given result of a search operation.
func (self *ProposalCache) RegisterSearchResult(source gitdomain.LocalBranchName, proposals []Proposal) {
	self.removeSearchResult(source)
	self.results = append(self.results, searchResult{
		source:    source,
		Proposals: proposals,
	})
}

func (self *ProposalCache) removeLookupResult(source, target gitdomain.LocalBranchName) {
	self.results = slices.DeleteFunc(self.results, func(result Result) bool {
		switch result := result.(type) {
		case lookupResult:
			return result.source == source && result.target == target
		}
		return false
	})
}

func (self *ProposalCache) removeSearchResult(source gitdomain.LocalBranchName) {
	self.results = slices.DeleteFunc(self.results, func(result Result) bool {
		switch result := result.(type) {
		case searchResult:
			return result.source == source
		}
		return false
	})
}
