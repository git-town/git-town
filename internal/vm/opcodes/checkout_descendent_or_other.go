package opcodes

import (
	"errors"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutAncestorOrOther checks out the first available ancestor branch of the current branch
// or any other branch available in the current worktree.
type CheckoutDescendentOrOther struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutDescendentOrOther) Run(args shared.RunArgs) error {
	availableBranches, err := args.Git.BranchesAvailableInCurrentWorktree(args.Backend)
	if err != nil {
		return err
	}
	availableBranches = availableBranches.Remove(self.Branch)
	descendents := args.Config.Value.NormalConfig.Lineage.Descendants(self.Branch, args.Config.Value.NormalConfig.Order)
	ancestors := args.Config.Value.NormalConfig.Lineage.Ancestors(self.Branch)
	slices.Reverse(ancestors)
	branches := availableBranches.Hoist(ancestors...).Hoist(descendents...)
	for _, branch := range branches {
		args.PrependOpcodes(&CheckoutIfNeeded{Branch: branch})
		return nil
	}
	return errors.New(messages.BranchNotAvailable)
}
