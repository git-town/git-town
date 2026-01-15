package github

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// type checks
var (
	mockAPIConnector MockConnector
	_                forgedomain.Connector = &mockAPIConnector
)

// MockConnector provides access to the Bitbucket Cloud API while caching proposal information.
type MockConnector struct {
	WebConnector
	Proposals     mockproposals.MockProposals
	ProposalsPath mockproposals.MockProposalPath
	cache         forgedomain.APICache
	log           print.Logger
}

// ============================================================================
// find proposals
// ============================================================================

var _ forgedomain.ProposalFinder = &mockAPIConnector // type check

func (self *MockConnector) FindProposal(source, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	if cachedProposal, has := self.cache.Lookup(source, target); has {
		return cachedProposal, nil
	}
	self.log.Start(messages.APIProposalFindStart, source, target)
	data, has := self.Proposals.FindBySourceAndTarget(source, target).Get()
	if !has {
		self.log.Success("none")
		return None[forgedomain.Proposal](), nil
	}
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+data.Number.String()), data.Title))
	proposal := forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeGithub}
	self.cache.RegisterLookupResult(source, target, Some(proposal))
	return Some(proposal), nil
}

// ============================================================================
// search proposals
// ============================================================================

var _ forgedomain.ProposalSearcher = &mockAPIConnector // type check

func (self *MockConnector) SearchProposals(source gitdomain.LocalBranchName) ([]forgedomain.Proposal, error) {
	if cachedSearchResult, has := self.cache.LookupSearch(source).Get(); has {
		return cachedSearchResult, nil
	}
	self.log.Start(messages.APIParentBranchLookupStart, source.String())
	result := []forgedomain.Proposal{}
	for _, data := range self.Proposals.FindBySource(source) {
		self.log.Success(data.Target.String())
		result = append(result, forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeGithub})
	}
	if len(result) == 0 {
		self.log.Success("none")
	}
	return result, nil
}

// ============================================================================
// update proposal body
// ============================================================================

var _ forgedomain.ProposalBodyUpdater = &mockAPIConnector // type check

func (self *MockConnector) UpdateProposalBody(proposalData forgedomain.ProposalInterface, newBody gitdomain.ProposalBody) error {
	self.cache.Clear()
	self.log.Start(messages.APIProposalUpdateBody, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()))
	proposal, hasProposal := self.Proposals.FindByID(proposalData.Data().Number).Get()
	if !hasProposal {
		return fmt.Errorf("proposal with id %d not found", proposalData.Data().Number)
	}
	proposal.Body = Some(newBody)
	self.Proposals.Update(proposal)
	mockproposals.Save(self.ProposalsPath, self.Proposals)
	self.log.Finished(nil)
	return nil
}

// ============================================================================
// update proposal target
// ============================================================================

var _ forgedomain.ProposalTargetUpdater = &mockAPIConnector // type check

func (self *MockConnector) UpdateProposalTarget(proposalData forgedomain.ProposalInterface, target gitdomain.LocalBranchName) error {
	self.cache.Clear()
	self.log.Start(messages.APIUpdateProposalTarget, colors.BoldGreen().Styled("#"+proposalData.Data().Number.String()), colors.BoldCyan().Styled(target.String()))
	proposal, hasProposal := self.Proposals.FindByID(proposalData.Data().Number).Get()
	if !hasProposal {
		return fmt.Errorf("proposal with id %d not found", proposalData.Data().Number)
	}
	proposal.Target = target
	self.Proposals.Update(proposal)
	mockproposals.Save(self.ProposalsPath, self.Proposals)
	self.log.Finished(nil)
	return nil
}
