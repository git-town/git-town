package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// UpdateProposalBase updates the target of the proposal with the given number at the code hosting platform.
type UpdateProposalBase struct {
	NewTarget               gitdomain.LocalBranchName
	OldTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateProposalBase) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *UpdateProposalBase) Run(args shared.RunArgs) error {
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		return connector.UpdateProposalBase(self.ProposalNumber, self.NewTarget, args.FinalMessages)
	}
	return hostingdomain.UnsupportedServiceError()
}

func (self *UpdateProposalBase) ShouldUndoOnError() bool {
	return true
}

func (self *UpdateProposalBase) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&UpdateProposalBase{
			NewTarget:      self.OldTarget,
			ProposalNumber: self.ProposalNumber,
		},
	}
}
