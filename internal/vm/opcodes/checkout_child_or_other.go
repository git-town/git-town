package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutChildOrOther checks out the first available child branch of the current branch,
// or any other branch available in the current worktree.
type CheckoutChildOrOther struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutChildOrOther) Run(args shared.RunArgs) error {
	availableBranches, err := args.Git.BranchesAvailableInCurrentWorktree(args.Backend)
	if err != nil {
		return err
	}
	availableBranches = availableBranches.Remove(self.Branch)
	descendents := args.Config.Value.NormalConfig.Lineage.Descendants(self.Branch, args.Config.Value.NormalConfig.Order)
	ancestors := args.Config.Value.NormalConfig.Lineage.Ancestors(self.Branch)
	branches := availableBranches.Hoist(descendents...).Hoist(ancestors...)
	for _, branch := range branches {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: branch})
		return nil
	}
	return errors.New("no branch to switch to available")
}
