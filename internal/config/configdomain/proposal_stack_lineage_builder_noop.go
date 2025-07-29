package configdomain

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// noopProposalStackLineageBuilder used by the stack lineage builder
// when the forge connector does not support finding proposals.
type noopProposalStackLineageBuilder struct{}

func (self *noopProposalStackLineageBuilder) AddBranch(childBranch gitdomain.LocalBranchName, parentBranch Option[gitdomain.LocalBranchName]) (ProposalStackLineageBuilder, error) {
	return self, nil
}

func (self *noopProposalStackLineageBuilder) Build(cfgs ...configureProposalStackLineage) Option[string] {
	return None[string]()
}

func (self *noopProposalStackLineageBuilder) GetProposal(branch gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	return None[forgedomain.ProposalData]()
}
