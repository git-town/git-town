package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// UpdateProposalTarget updates the target of the proposal with the given number at the code hosting platform.
type UpdateProposalTarget struct {
	NewTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateProposalTarget) CreateAutomaticUndoError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}

func (self *UpdateProposalTarget) Run(args shared.RunArgs) error {
	if connector, hasConnector := args.Connector.Get(); hasConnector {
		return connector.UpdateProposalTarget(self.ProposalNumber, self.NewTarget)
	}
	return hostingdomain.UnsupportedServiceError()
}

func (self *UpdateProposalTarget) ShouldAutomaticallyUndoOnError() bool {
	return true
}
