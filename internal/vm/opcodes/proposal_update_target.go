package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// ProposalUpdateTarget updates the target of the proposal with the given number at the forge.
type ProposalUpdateTarget struct {
	NewBranch               gitdomain.LocalBranchName
	OldBranch               gitdomain.LocalBranchName
	Proposal                forgedomain.Proposal
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateTarget) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.Proposal.Data.Data().Number)
}

func (self *ProposalUpdateTarget) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	updateProposalTarget, canUpdateProposalTarget := connector.UpdateProposalTargetFn().Get()
	if !canUpdateProposalTarget {
		return forgedomain.UnsupportedServiceError()
	}
	return updateProposalTarget(self.Proposal.Data, self.NewBranch, args.FinalMessages)
}

func (self *ProposalUpdateTarget) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateTarget) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateTarget{
			NewBranch: self.OldBranch,
			OldBranch: self.NewBranch,
			Proposal:  self.Proposal,
		},
	}
}
