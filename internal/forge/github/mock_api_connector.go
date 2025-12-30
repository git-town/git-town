package github

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	mockAPIConnector MockAPIConnector
	_                forgedomain.Connector = &mockAPIConnector
)

// MockAPIConnector provides access to the Bitbucket Cloud API while caching proposal information.
type MockAPIConnector struct {
	WebConnector
	Proposals mockproposals.MockProposals
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &mockAPIConnector // type check

func (self *MockAPIConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	data, has := self.Proposals.FindBySourceAndTarget(source, target).Get()
	if !has {
		return None[forgedomain.Proposal](), nil
	}
	return Some(forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeGitHub}), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &mockAPIConnector // type check

func (self *MockAPIConnector) SearchProposals(source gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	result := []forgedomain.Proposal{}
	for _, data := range self.Proposals.Search(source) {
		result = append(result, forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeGitHub})
	}
	return result, nil
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &mockAPIConnector // type check

func (self *MockAPIConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody gitdomain.ProposalBody) error {
	fmt.Println("1111111111111111111111111111111111111111111111111111111111111111111111111")
	proposal, hasProposal := self.Proposals.FindByID(proposalData.Data().Number).Get()
	if !hasProposal {
		return fmt.Errorf("proposal with id %d not found", proposalData.Data().Number)
	}
	proposal.Body = Some(newBody)
	return nil
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &mockAPIConnector // type check

func (self *MockAPIConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	proposal, hasProposal := self.Proposals.FindByID(proposalData.Data().Number).Get()
	if !hasProposal {
		return fmt.Errorf("proposal with id %d not found", proposalData.Data().Number)
	}
	proposal.Target = target
	return nil
}
