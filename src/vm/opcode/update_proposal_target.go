package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// UpdateProposalTarget updates the target of the proposal with the given number at the code hosting service.
type UpdateProposalTarget struct {
	ProposalNumber int
	NewTarget      domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *UpdateProposalTarget) Run(args shared.RunArgs) error {
	return args.Connector.UpdateProposalTarget(self.ProposalNumber, self.NewTarget)
}

func (self *UpdateProposalTarget) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (self *UpdateProposalTarget) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, self.ProposalNumber)
}
