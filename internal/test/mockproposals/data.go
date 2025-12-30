package mockproposals

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type MockProposals struct {
	FilePath  string // path of the JSON file containing the proposal data on disk
	Proposals []forgedomain.ProposalData
}

func (self *MockProposals) FindByID(id int) OptionalMutable[forgedomain.ProposalData] {
	for p, proposal := range self.Proposals {
		if proposal.Number == id {
			return MutableSome(&self.Proposals[p])
		}
	}
	return MutableNone[forgedomain.ProposalData]()
}

func (self *MockProposals) FindBySourceAndTarget(source, target gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	for p, proposal := range self.Proposals {
		if proposal.Source == source && proposal.Target == target {
			return Some(self.Proposals[p])
		}
	}
	return None[forgedomain.ProposalData]()
}

// Save stores the changes made to the given proposal to disk.
func (self *MockProposals) Save(proposal forgedomain.ProposalData) {
	self.Update(proposal)
	Save(self.FilePath, self.Proposals)
}

func (self *MockProposals) Search(source gitdomain.LocalBranchName) []forgedomain.ProposalData {
	result := []forgedomain.ProposalData{}
	for _, proposal := range self.Proposals {
		if proposal.Source == source {
			result = append(result, proposal)
		}
	}
	return result
}

func (self *MockProposals) Update(proposal forgedomain.ProposalData) {
	for p := range self.Proposals {
		if self.Proposals[p].Number == proposal.Number {
			self.Proposals[p] = proposal
		}
	}
}
