package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CheckoutStep checks out a new branch.
type CheckoutStep struct {
	EmptyStep
	Branch         domain.LocalBranchName
	previousBranch domain.LocalBranchName
}

func (step *CheckoutStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CheckoutStep{Branch: step.previousBranch}}, nil
}

func (step *CheckoutStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousBranch, err = run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if step.previousBranch != step.Branch {
		err := run.Frontend.CheckoutBranch(step.Branch)
		return err
	}
	return nil
}
