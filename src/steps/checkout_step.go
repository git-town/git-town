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
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == step.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(step.Branch)
}
