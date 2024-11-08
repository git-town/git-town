package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ProposalUpdateBase updates the target of the proposal with the given number at the code hosting platform.
type ProposalUpdateBase struct {
	NewTarget               gitdomain.LocalBranchName
	OldTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateBase) AutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *ProposalUpdateBase) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return hostingdomain.UnsupportedServiceError()
	}
	updateProposalTarget, canUpdateProposalTarget := connector.UpdateProposalTargetFn().Get()
	if !canUpdateProposalTarget {
		return hostingdomain.UnsupportedServiceError()
	}
	return updateProposalTarget(self.ProposalNumber, self.NewTarget, args.FinalMessages)
}

func (self *ProposalUpdateBase) ShouldUndoOnError() bool {
	return true
}

func (self *ProposalUpdateBase) UndoExternalChangesProgram() []shared.Opcode {
	return []shared.Opcode{
		&ProposalUpdateBase{
			NewTarget:      self.OldTarget,
			OldTarget:      self.NewTarget,
			ProposalNumber: self.ProposalNumber,
		},
	}
}
