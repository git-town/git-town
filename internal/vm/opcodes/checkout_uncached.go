package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// CheckoutIfNeeded checks out a new branch.
type CheckoutUncached struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutUncached) Run(args shared.RunArgs) error {
	_ = args.Git.CheckoutBranchUncached(args.Backend, self.Branch, false)
	return nil
}
