package forgedomain

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCache caches proposals by source and target branches.
type ProposalCache struct {
	// the cached proposals
	//
	// A simple list is fine here despite O(n) lookup time because the number of proposals is expected to be very small.
	proposals []Proposal
}

// BySource provides the cached proposals for the given source branch.
func (self *ProposalCache) BySource(source gitdomain.LocalBranchName) []Proposal {
	result := []Proposal{}
	for _, proposal := range self.proposals {
		if proposal.Data.Data().Source == source {
			result = append(result, proposal)
		}
	}
	return result
}

// BySourceTarget provides the cached proposal for the given source and target branch.
func (self *ProposalCache) BySourceTarget(source, target gitdomain.LocalBranchName) Option[Proposal] {
	for _, proposal := range self.proposals {
		proposalData := proposal.Data.Data()
		if proposalData.Source == source && proposalData.Target == target {
			return Some(proposal)
		}
	}
	return None[Proposal]()
}

func (self *ProposalCache) Delete(proposalNumber int) {
	self.proposals = slices.DeleteFunc(self.proposals, func(proposal Proposal) bool {
		return proposal.Data.Data().Number == proposalNumber
	})
}

// Set caches the given proposal.
func (self *ProposalCache) Set(proposal Proposal) {
	self.proposals = append(self.proposals, proposal)
}

func (self *ProposalCache) SetMany(proposals []Proposal) {
	self.proposals = append(self.proposals, proposals...)
}

func (self *ProposalCache) SetOption(proposalOpt Option[Proposal]) {
	if proposal, hasProposal := proposalOpt.Get(); hasProposal {
		self.Set(proposal)
	}
}
