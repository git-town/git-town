package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// Checkout checks out a new branch.
type Checkout struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *Checkout) Run(args shared.RunArgs) error {
	return args.Git.CheckoutBranch(args.Frontend, self.Branch, false)
}
