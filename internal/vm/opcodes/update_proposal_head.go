package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// UpdateProposalHead updates the head of the proposal with the given number at the code hosting platform.
type UpdateProposalHead struct {
	NewTarget               gitdomain.LocalBranchName
	OldTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateProposalHead) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *UpdateProposalHead) Run(args shared.RunArgs) error {
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		return connector.UpdateProposalHead(self.ProposalNumber, self.NewTarget, args.FinalMessages)
	}
	return hostingdomain.UnsupportedServiceError()
}

func (self *UpdateProposalHead) ShouldUndoOnError() bool {
	return true
}

func (self *UpdateProposalHead) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&UpdateProposalHead{
			NewTarget:      self.OldTarget,
			ProposalNumber: self.ProposalNumber,
		},
	}
}
