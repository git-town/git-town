package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/vm/shared"
)

// ProposalUpdateSource updates the source branch of the proposal with the given number at the forge.
type ProposalUpdateSource struct {
	NewBranch               gitdomain.LocalBranchName
	OldBranch               gitdomain.LocalBranchName
	Proposal                forgedomain.SerializableProposal
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateSource) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.Proposal.Data.GetNumber())
}

func (self *ProposalUpdateSource) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	updateProposalSource, canUpdateProposalSource := connector.UpdateProposalSourceFn().Get()
	if !canUpdateProposalSource {
		return errors.New(messages.ProposalSourceCannotUpdate)
	}
	return updateProposalSource(self.Proposal, self.NewBranch, args.FinalMessages)
}

func (self *ProposalUpdateSource) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateSource) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateSource{
			NewBranch: self.OldBranch,
			OldBranch: self.NewBranch,
			Proposal:  self.Proposal,
		},
	}
}
