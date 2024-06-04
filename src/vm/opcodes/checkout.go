package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *Checkout) Run(args shared.RunArgs) error {
	existingBranch, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	if existingBranch == self.Branch {
		return nil
	}
	return args.Git.CheckoutBranch(args.Frontend, self.Branch, false)
}
