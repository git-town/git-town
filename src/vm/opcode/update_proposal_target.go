package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type UpdateProposalTarget struct {
	ProposalNumber int
	NewTarget      domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *UpdateProposalTarget) Run(args shared.RunArgs) error {
	return args.Connector.UpdateProposalTarget(step.ProposalNumber, step.NewTarget)
}

func (step *UpdateProposalTarget) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (step *UpdateProposalTarget) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, step.ProposalNumber)
}
