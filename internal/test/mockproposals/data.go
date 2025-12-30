package mockproposals

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type MockProposals []forgedomain.ProposalData

func (self MockProposals) FindByID(id int) OptionalMutable[forgedomain.ProposalData] {
	for _, proposal := range self {
		if proposal.Number == id {
			return MutableSome(&proposal)
		}
	}
	return MutableNone[forgedomain.ProposalData]()
}

func (self MockProposals) FindBySourceAndTarget(source, target gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	for _, proposal := range self {
		if proposal.Source == source && proposal.Target == target {
			return Some(proposal)
		}
	}
	return None[forgedomain.ProposalData]()
}

func (self MockProposals) Search(source gitdomain.LocalBranchName) []forgedomain.ProposalData {
	result := []forgedomain.ProposalData{}
	for _, proposal := range self {
		if proposal.Source == source {
			result = append(result, proposal)
		}
	}
	return result
}
