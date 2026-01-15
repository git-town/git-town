package gitea

import (
	"fmt"
	"strconv"

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
	self.log.Log(fmt.Sprintf("%s (%s)", colors.BoldGreen().Styled("#"+strconv.Itoa(data.Number)), data.Title))
	return Some(forgedomain.Proposal{Data: data, ForgeType: forgedomain.ForgeTypeGithub}), nil
}
