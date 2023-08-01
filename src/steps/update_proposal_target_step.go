package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type UpdateProposalTargetStep struct {
	ProposalNumber int
	NewTarget      string
	ExistingTarget string
	EmptyStep
}

func (step *UpdateProposalTargetStep) Run(_ *git.ProdRunner, connector hosting.Connector) error {
	return connector.UpdateProposalTarget(step.ProposalNumber, step.NewTarget)
}

func (step *UpdateProposalTargetStep) CreateAbortStep() Step {
	return &step.EmptyStep
}

func (step *UpdateProposalTargetStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&UpdateProposalTargetStep{
		ProposalNumber: step.ProposalNumber,
		NewTarget:      step.ExistingTarget,
		ExistingTarget: step.NewTarget,
	}}, nil
}

func (step *UpdateProposalTargetStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (step *UpdateProposalTargetStep) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ProposalTargetBranchUpdateProblem, step.ProposalNumber)
}
