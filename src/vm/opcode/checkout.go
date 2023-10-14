package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch domain.LocalBranchName
	BaseOpcode
}

func (step *Checkout) Run(args RunArgs) error {
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == step.Branch {
		return nil
	}
	return args.Runner.Frontend.CheckoutBranch(step.Branch)
}
