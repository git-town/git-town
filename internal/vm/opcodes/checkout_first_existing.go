package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
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
		args.PrependOpcodes(&Checkout{Branch: existingBranch})
	} else {
		args.PrependOpcodes(&Checkout{Branch: self.MainBranch})
	}
	return nil
}
