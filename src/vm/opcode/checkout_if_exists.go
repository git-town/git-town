package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CheckoutIfExists does the same as CheckoutStep
// but only if that branch actually exists.
type CheckoutIfExists struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (step *CheckoutIfExists) Run(args shared.RunArgs) error {
	if !args.Runner.Backend.HasLocalBranch(step.Branch) {
		return nil
	}
	return (&Checkout{Branch: step.Branch}).Run(args)
}
