package step

import "github.com/git-town/git-town/v9/src/domain"

// Checkout checks out a new branch.
type Checkout struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *Checkout) Run(args RunArgs) error {
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if existingBranch == step.Branch {
		return nil
	}
	err = args.Runner.Frontend.CheckoutBranch(step.Branch)
	return err
}
