package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// CheckoutStep checks out a new branch.
type CheckoutStep struct {
	Branch         domain.LocalBranchName
	previousBranch domain.LocalBranchName `exhaustruct:"optional"`
	EmptyStep
}

func (step *CheckoutStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CheckoutStep{Branch: step.previousBranch}}, nil
}

func (step *CheckoutStep) Run(args RunArgs) error {
	var err error
	step.previousBranch, err = args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if step.previousBranch != step.Branch {
		err := args.Runner.Frontend.CheckoutBranch(step.Branch)
		return err
	}
	return nil
}
