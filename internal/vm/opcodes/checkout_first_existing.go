package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutFirstExisting checks out the first existing branch,
// but only if that branch actually exists.
type CheckoutFirstExisting struct {
	Branches   gitdomain.LocalBranchNames
	MainBranch gitdomain.LocalBranchName
}

func (self *CheckoutFirstExisting) Run(args shared.RunArgs) error {
	if existingBranch, hasExistingBranch := args.Git.FirstExistingBranch(args.Backend, self.Branches...).Get(); hasExistingBranch {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: existingBranch})
	} else {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: self.MainBranch})
	}
	return nil
}
