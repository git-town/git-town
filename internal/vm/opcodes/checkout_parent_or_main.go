package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutChildOrOther checks out the first available child branch of the current branch,
// or any other branch available in the current worktree.
type CheckoutChildOrOther struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutChildOrOther) Run(args shared.RunArgs) error {
	// first try to checkout the first available child branch
	children := args.Config.Value.NormalConfig.Lineage.Children(self.Branch, args.Config.Value.NormalConfig.Order)
	availableBranches, err := args.Git.BranchesAvailableInCurrentWorktree(args.Backend)
	for _, child := range children {
		if args.Git.BranchExists(args.Backend, child) {
			args.PrependOpcodes(&CheckoutIfNeeded{Branch: child})
			return nil
		}
	}
	if len(children) > 0 {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: children[0]})
		return nil
	}
	ancestors := args.Config.Value.NormalConfig.Lineage.Ancestors(self.Branch)
	if len(ancestors) > 0 {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: ancestors[0]})
	}
	args.PrependOpcodes(&CheckoutIfNeeded{Branch: parent})
	return nil
}
