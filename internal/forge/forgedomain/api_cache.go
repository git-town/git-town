package forgedomain

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// APICache caches results of API calls to a forge.
type APICache struct {
	results []Result // the cached items
}

type Result interface {
	isAPIResult() // type safety to allow only lookupResult and searchResult
}

// the result of a find operation
type lookupResult struct {
	proposal Option[Proposal]
	source   gitdomain.LocalBranchName
	target   gitdomain.LocalBranchName
}

func (lr lookupResult) isAPIResult() {}

// the result of a search operation
type searchResult struct {
	proposals []Proposal
	source    gitdomain.LocalBranchName
}

func (sr searchResult) isAPIResult() {}

// Clear removes all cached results.
func (self *APICache) Clear() {
	self.results = []Result{}
}

// Lookup provides what this cache knows about the proposal for the given source and target branch.
// If it has a cached proposal, returns (Some, true).
// If it knowns that there is no proposal, returns (None, true).
// If it has never heard about the given source and target branch, returns (None, false).
func (self *APICache) Lookup(source, target gitdomain.LocalBranchName) (proposal Option[Proposal], knows bool) {
	for _, result := range self.results {
		switch result := result.(type) {
		case lookupResult:
			if result.source == source && result.target == target {
				return result.proposal, true
			}
		case searchResult:
			if result.source == source {
				for _, proposal := range result.proposals {
					if proposal.Data.Data().Target == target {
						return Some(proposal), true
					}
				}
				// here we know that there was a proposal search for the source branch,
				// and the target branch was not in the result --> we know there is no proposal
				return None[Proposal](), true
			}
		default:
			panic(fmt.Sprintf("unknown result type: %T", result))
		}
	}
	// we didn't run across any API results for the source branch
	return None[Proposal](), false
}

// LookupSearch provides the cached search result for the given source branch.
func (self *APICache) LookupSearch(source gitdomain.LocalBranchName) (proposals []Proposal, knows bool) {
	for _, result := range self.results {
		if searchResult, ok := result.(searchResult); ok {
			if searchResult.source == source {
				return searchResult.proposals, true
			}
		}
	}
	return []Proposal{}, false
}

// SaveLookupResult registers the given result of a lookup operation.
func (self *APICache) RegisterLookupResult(source, target gitdomain.LocalBranchName, proposal Option[Proposal]) {
	self.removeLookupResult(source, target)
	self.results = append(self.results, lookupResult{
		proposal: proposal,
		source:   source,
		target:   target,
	})
}

// RegisterSearchResult registers the given result of a search operation.
func (self *APICache) RegisterSearchResult(source gitdomain.LocalBranchName, proposals []Proposal) {
	self.removeSearchResult(source)
	self.results = append(self.results, searchResult{
		proposals: proposals,
		source:    source,
	})
}

func (self *APICache) removeLookupResult(source, target gitdomain.LocalBranchName) {
	self.results = slices.DeleteFunc(self.results, func(result Result) bool {
		if result, ok := result.(lookupResult); ok {
			return result.source == source && result.target == target
		}
		return false
	})
}

func (self *APICache) removeSearchResult(source gitdomain.LocalBranchName) {
	self.results = slices.DeleteFunc(self.results, func(result Result) bool {
		if result, ok := result.(searchResult); ok {
			return result.source == source
		}
		return false
	})
}
