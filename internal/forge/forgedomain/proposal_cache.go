package forgedomain

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCache caches known proposals and knowledge about non-existing proposals.
type ProposalCache struct {
	// the cached items
	//
	// A simple list is fine here despite O(n) lookup time because the number of items is expected to be very small.
	items []proposalCacheItem
}

// proposalCacheItem encodes knowledge whether a proposal for a given source and target branch exists.
type proposalCacheItem struct {
	Source   gitdomain.LocalBranchName
	Target   gitdomain.LocalBranchName
	Proposal Option[Proposal] // Some if a proposal exists, None if it doesn't
}

// Lookup provides what this cache knows about the proposal for the given source and target branch.
// If it has a cached proposal, returns (Some, true).
// If it knowns that there is no proposal, returns (None, true).
// If it has never heard about the given source and target branch, returns (None, false).
func (self *ProposalCache) Lookup(source, target gitdomain.LocalBranchName) (proposal Option[Proposal], knows bool) {
	for _, item := range self.items {
		if item.Source == source && item.Target == target {
			return item.Proposal, true
		}
	}
	return None[Proposal](), false
}

// DeleteBySourceTarget removes the cached proposal for the given source and target branch.
func (self *ProposalCache) DeleteBySourceTarget(source, target gitdomain.LocalBranchName) {
	self.items = slices.DeleteFunc(self.items, func(item proposalCacheItem) bool {
		return item.Source == source && item.Target == target
	})
}

// DeleteByNumber removes the cached proposal with the given number.
func (self *ProposalCache) DeleteByNumber(proposalNumber int) {
	self.items = slices.DeleteFunc(self.items, func(item proposalCacheItem) bool {
		if proposal, hasProposal := item.Proposal.Get(); hasProposal {
			return proposal.Data.Data().Number == proposalNumber
		}
		return false
	})
}

func (self *ProposalCache) Delete(proposal ProposalInterface) {
	proposalData := proposal.Data()
	self.items = slices.DeleteFunc(self.items, func(item proposalCacheItem) bool {
		return item.Source == proposalData.Source && item.Target == proposalData.Target
	})
}

// Set registers knowledge about the proposal for the given source and target branch.
// If you provide None for the proposal, it stores knowledge that there is no proposal for the given source and target branch.
func (self *ProposalCache) Set(source, target gitdomain.LocalBranchName, proposal Option[Proposal]) {
	self.DeleteBySourceTarget(source, target)
	self.items = append(self.items, proposalCacheItem{
		Source:   source,
		Target:   target,
		Proposal: proposal,
	})
}

// SetMany caches knowledge about many proposals.
func (self *ProposalCache) SetMany(proposals []Proposal) {
	for _, proposal := range proposals {
		proposalData := proposal.Data.Data()
		self.Set(proposalData.Source, proposalData.Target, Some(proposal))
	}
}

// Search provides the cached proposals for the given source branch.
// If it is known that the source branch has no proposals, return (empty, true).
// If it isn't known whether this branch has proposals, return (empty, false).
func (self *ProposalCache) Search(source gitdomain.LocalBranchName) (proposals []Proposal, knows bool) {
	for _, item := range self.items {
		if item.Source == source {
			knows = true
			if proposal, hasProposal := item.Proposal.Get(); hasProposal {
				proposals = append(proposals, proposal)
			}
		}
	}
	return
}
