package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type UpdateProposalTargetStep struct {
	ProposalNumber int
	NewTarget      string
	ExistingTarget string
	EmptyStep
}

func (step *UpdateProposalTargetStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return connector.UpdateProposalTarget(step.ProposalNumber, step.NewTarget)
}

func (step *UpdateProposalTargetStep) CreateAbortStep() Step {
	return &step.EmptyStep
}

func (step *UpdateProposalTargetStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &UpdateProposalTargetStep{
		ProposalNumber: step.ProposalNumber,
		NewTarget:      step.ExistingTarget,
		ExistingTarget: step.NewTarget,
	}, nil
}

func (step *UpdateProposalTargetStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}

func (step *UpdateProposalTargetStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("cannot update the target branch of proposal %d via the API", step.ProposalNumber)
}
