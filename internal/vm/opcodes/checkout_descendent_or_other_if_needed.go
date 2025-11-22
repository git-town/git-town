package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutDescendentOrOtherIfNeeded checks out the first available ancestor branch of the current branch
// or any other branch available in the current worktree
// if the current branch is the same as the given branch.
type CheckoutDescendentOrOtherIfNeeded struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutDescendentOrOtherIfNeeded) Run(args shared.RunArgs) error {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
		if currentBranch != self.Branch {
			return nil
		}
	}
	args.PrependOpcodes(&CheckoutDescendentOrOther{Branch: self.Branch})
	return nil
}
