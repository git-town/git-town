package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ProposalUpdateHead updates the head of the proposal with the given number at the code hosting platform.
type ProposalUpdateHead struct {
	NewTarget               gitdomain.LocalBranchName
	OldTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateHead) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *ProposalUpdateHead) Run(args shared.RunArgs) error {
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		return connector.UpdateProposalHead(self.ProposalNumber, self.NewTarget, args.FinalMessages)
	}
	return hostingdomain.UnsupportedServiceError()
}

func (self *ProposalUpdateHead) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateHead) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateHead{
			NewTarget:      self.OldTarget,
			OldTarget:      self.NewTarget,
			ProposalNumber: self.ProposalNumber,
		},
	}
}
