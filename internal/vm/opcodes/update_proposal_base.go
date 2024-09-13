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
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateProposalBase) CreateAutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *UpdateProposalBase) Run(args shared.RunArgs) error {
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		return connector.UpdateProposalBase(self.ProposalNumber, self.NewTarget)
	}
	return hostingdomain.UnsupportedServiceError()
}

func (self *UpdateProposalBase) ShouldAutomaticallyUndoOnError() bool {
	return true
}
