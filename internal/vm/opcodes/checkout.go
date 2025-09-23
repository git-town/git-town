package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// Checkout checks out the given existing branch.
type Checkout struct {
	Branch gitdomain.LocalBranchName
}

func (self *Checkout) Run(args shared.RunArgs) error {
	return args.Git.CheckoutBranch(args.Frontend, self.Branch, false)
}
