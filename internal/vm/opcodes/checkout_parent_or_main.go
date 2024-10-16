package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// CheckoutParentOrMain checks out the parent branch of the current branch,
// or the main branch if the current branch has no parent.
type CheckoutParentOrMain struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CheckoutParentOrMain) Run(args shared.RunArgs) error {
	parent := args.Config.Config.Lineage.Parent(self.Branch).GetOrElse(args.Config.Config.MainBranch)
	args.PrependOpcodes(&CheckoutIfNeeded{Branch: parent})
	return nil
}
