package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// CheckoutIfExistsStep does the same as CheckoutStep
// but only if that branch actually exists.
type CheckoutIfExistsStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *CheckoutIfExistsStep) Run(args RunArgs) error {
	if !args.Runner.Backend.HasLocalBranch(step.Branch) {
		return nil
	}
	previousBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if previousBranch == step.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(step.Branch)
}
