package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoOpStep
	BranchName string

	previousBranchName string
}

// CreateUndoStep returns the undo step for this step.
func (step *CheckoutBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &CheckoutBranchStep{BranchName: step.previousBranchName}, nil
}

// Run executes this step.
func (step *CheckoutBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	step.previousBranchName = git.GetCurrentBranchName()
	if step.previousBranchName != step.BranchName {
		err := repo.Logging.CheckoutBranch(step.BranchName)
		if err == nil {
			git.UpdateCurrentBranchCache(step.BranchName)
		}
		return err
	}
	return nil
}
