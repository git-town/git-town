package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CheckoutIfNeeded checks out a new branch.
type CheckoutUncached struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutUncached) Run(args shared.RunArgs) error {
	return args.Git.CheckoutBranchUncached(args.Frontend, self.Branch, false)
}
