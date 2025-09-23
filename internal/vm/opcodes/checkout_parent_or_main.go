package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutParentOrMain checks out the parent branch of the current branch,
// or the main branch if the current branch has no parent.
type CheckoutParentOrMain struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutParentOrMain) Run(args shared.RunArgs) error {
	parent := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).GetOr(args.Config.Value.ValidatedConfigData.MainBranch)
	args.PrependOpcodes(&CheckoutIfNeeded{Branch: parent})
	return nil
}
