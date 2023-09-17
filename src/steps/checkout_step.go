package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CheckoutStep checks out a new branch.
type CheckoutStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *CheckoutStep) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if step.Branch == currentBranch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(step.Branch)
}
