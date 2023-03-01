package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CheckoutStep checks out a new branch.
type CheckoutStep struct {
	EmptyStep
	Branch         string
	previousBranch string
}

func (step *CheckoutStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &CheckoutStep{Branch: step.previousBranch}, nil
}

func (step *CheckoutStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
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
