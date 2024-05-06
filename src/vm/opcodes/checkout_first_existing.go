package opcodes

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CheckoutIfExists does the same as Checkout
// but only if that branch actually exists.
type CheckoutFirstExisting struct {
	Branches                gitdomain.LocalBranchNames
	MainBranch              gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutFirstExisting) Run(args shared.RunArgs) error {
	existingBranch := args.Backend.FirstExistingBranch(self.Branches, self.MainBranch)
	return (&Checkout{Branch: existingBranch}).Run(args)
}
