package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutUncached checks out a new branch ignoring the current-branch cache..
type CheckoutUncached struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutUncached) Run(args shared.RunArgs) error {
	_ = args.Git.CheckoutBranchUncached(args.Backend, self.Branch, false)
	return nil
}
