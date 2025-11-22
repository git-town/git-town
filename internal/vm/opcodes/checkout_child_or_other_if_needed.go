package opcodes

import (
	"errors"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CheckoutAncestorOrOtherIfNeeded checks out the first available ancestor branch of the current branch
// or any other branch available in the current worktree
// if the current branch is the same as the given branch.
type CheckoutAncestorOrOtherIfNeeded struct {
	Branch gitdomain.LocalBranchName
}

func (self *CheckoutAncestorOrOtherIfNeeded) Run(args shared.RunArgs) error {
	currentBranchOpt, err := args.Git.CurrentBranch(args.Backend)
	if err != nil {
		return err
	}
	if currentBranch, hasCurrentBranch := currentBranchOpt.Get(); hasCurrentBranch {
		if currentBranch != self.Branch {
			return nil
		}
	}
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
