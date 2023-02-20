package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CheckoutBranchStep checks out a new branch.
type CheckoutBranchStep struct {
	NoOpStep
	Branch         string
	previousBranch string
}

func (step *CheckoutBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &CheckoutBranchStep{Branch: step.previousBranch}, nil
}

func (step *CheckoutBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	var err error
	step.previousBranch, err = repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	if step.previousBranch != step.Branch {
		err := repo.Logging.CheckoutBranch(step.Branch)
		return err
	}
	return nil
}
