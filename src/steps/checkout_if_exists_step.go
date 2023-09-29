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
	return (&CheckoutStep{Branch: step.Branch}).Run(args)
}
