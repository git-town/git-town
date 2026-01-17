package mockproposals

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type MockProposals []forgedomain.ProposalData

func (self *MockProposals) DeleteByID(id forgedomain.ProposalNumber) {
	*self = slices.DeleteFunc(*self, func(proposal forgedomain.ProposalData) bool {
		return proposal.Number == id
	})
}

func (self *MockProposals) FindByID(id forgedomain.ProposalNumber) Option[forgedomain.ProposalData] {
	for _, proposal := range *self {
		if proposal.Number == id {
			return Some(proposal)
		}
	}
	return None[forgedomain.ProposalData]()
}

func (self *MockProposals) FindBySource(source gitdomain.LocalBranchName) []forgedomain.ProposalData {
	result := []forgedomain.ProposalData{}
	for _, proposal := range *self {
		if proposal.Source == source {
			result = append(result, proposal)
		}
	}
	return result
}

func (self *MockProposals) FindBySourceAndTarget(source, target gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	for _, proposal := range *self {
		if proposal.Source == source && proposal.Target == target {
			return Some(proposal)
		}
	}
	return None[forgedomain.ProposalData]()
}

func (self *MockProposals) Update(proposal forgedomain.ProposalData) {
	for p, prop := range *self {
		if prop.Number == proposal.Number {
			(*self)[p] = proposal
		}
	}
}
