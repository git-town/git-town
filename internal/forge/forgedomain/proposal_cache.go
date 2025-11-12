package forgedomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCache caches proposals by source and target branches.
type ProposalCache struct {
	// the cached proposals
	//
	// A simply list is fine here despite O(n) lookup time because the number of proposals is expected to be very small.
	proposals []Proposal
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

// BySource provides the cached proposal for the given source branch.
func (self *ProposalCache) BySource(source gitdomain.LocalBranchName) Option[Proposal] {
	for _, proposal := range self.proposals {
		if proposal.Data.Data().Source == source {
			return Some(proposal)
		}
	}
	return None[Proposal]()
}

// Set caches the given proposal.
func (self *ProposalCache) Set(proposal Proposal) {
	self.proposals = append(self.proposals, proposal)
}
