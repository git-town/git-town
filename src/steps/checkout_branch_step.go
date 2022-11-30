package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoOpStep
	BranchName         string
	previousBranchName string
}

func (step *CheckoutBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &CheckoutBranchStep{BranchName: step.previousBranchName}, nil
}

func (step *CheckoutBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	var err error
	step.previousBranchName, err = repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	if step.previousBranchName != step.BranchName {
		err := repo.Logging.CheckoutBranch(step.BranchName)
		return err
	}
	return nil
}
