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

func (op *UpdateProposalTarget) Run(args shared.RunArgs) error {
	return args.Connector.UpdateProposalTarget(op.ProposalNumber, op.NewTarget)
}

func (op *UpdateProposalTarget) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (op *UpdateProposalTarget) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, op.ProposalNumber)
}
