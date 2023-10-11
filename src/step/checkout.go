package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *Checkout) Run(args RunArgs) error {
	fmt.Println("2222222222222222 CHECKOUT BRANCH")
	fmt.Println("2222222222222222 BRANCH TO CHECK OUT", step.Branch)
	existingBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	fmt.Println("2222222222222222 EXISTING BRANCH", existingBranch)
	if existingBranch == step.Branch {
		return nil
	}
	err = args.Runner.Frontend.CheckoutBranch(step.Branch)
	fmt.Println("2222222222222222 EXISTING BRANCH AFTER", args.Runner.Backend.CurrentBranchCache)
	return err
}
