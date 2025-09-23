package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// ProposalUpdateSource updates the source branch of the proposal with the given number at the forge.
type ProposalUpdateSource struct {
	NewBranch gitdomain.LocalBranchName
	OldBranch gitdomain.LocalBranchName
	Proposal  forgedomain.Proposal
}

func (self *ProposalUpdateSource) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.Proposal.Data.Data().Number)
}

func (self *ProposalUpdateSource) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	proposalSourceUpdater, canUpdateProposalSource := connector.(forgedomain.ProposalSourceUpdater)
	if !canUpdateProposalSource {
		return errors.New(messages.ProposalSourceCannotUpdate)
	}
	return proposalSourceUpdater.UpdateProposalSource(self.Proposal.Data, self.NewBranch)
}

func (self *ProposalUpdateSource) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateSource) UndoExternalChanges() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateSource{
			NewBranch: self.OldBranch,
			OldBranch: self.NewBranch,
			Proposal:  self.Proposal,
		},
	}
}
