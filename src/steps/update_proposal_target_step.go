package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type UpdateProposalTargetStep struct {
	ProposalNumber int
	NewTarget      domain.LocalBranchName
	EmptyStep
}

func (step *UpdateProposalTargetStep) Run(args RunArgs) error {
	return args.Connector.UpdateProposalTarget(step.ProposalNumber, step.NewTarget)
}

func (step *UpdateProposalTargetStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (step *UpdateProposalTargetStep) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, step.ProposalNumber)
}
