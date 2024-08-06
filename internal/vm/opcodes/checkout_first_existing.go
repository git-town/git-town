package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutFirstExisting struct {
	Branches                gitdomain.LocalBranchNames
	MainBranch              gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutFirstExisting) Run(args shared.RunArgs) error {
	if existingBranch, hasExistingBranch := args.Git.FirstExistingBranch(args.Backend, self.Branches...).Get(); hasExistingBranch {
		return (&Checkout{Branch: existingBranch}).Run(args)
	}
	return (&Checkout{Branch: self.MainBranch}).Run(args)
}
