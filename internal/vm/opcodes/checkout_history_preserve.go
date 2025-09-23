package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// CheckoutHistoryPreserve does stuff.
type CheckoutHistoryPreserve struct {
	PreviousBranchCandidates []Option[gitdomain.LocalBranchName]
}

func (self *CheckoutHistoryPreserve) Run(args shared.RunArgs) error {
	cachedCurrentBranch, hasCachedCurrentBranch := args.Git.CurrentBranchCache.Get()
	if !hasCachedCurrentBranch {
		// the branch cache is not initialized --> there were no branch changes --> no need to restore the branch history
		return nil
	}
	cachedPreviousBranch, hasCachedPreviousBranch := args.Git.CurrentBranchCache.GetPrevious()
	// remove the current branch from the list of previous branch candidates because the current branch should never also be the previous branch
	candidates := slice.GetAll(self.PreviousBranchCandidates)
	candidatesWithoutCurrent := slice.Remove(candidates, cachedCurrentBranch)
	existingPreviousBranch, hasExistingPreviousBranch := args.Git.FirstExistingBranch(args.Backend, candidatesWithoutCurrent...).Get()
	if !hasExistingPreviousBranch || !hasCachedPreviousBranch || cachedPreviousBranch == existingPreviousBranch {
		return nil
	}
	// We	need to ignore errors here because failing to set the Git branch history
	// is not an error condition.
	// This operation can fail for a number of reasons like the previous branch being
	// checked out in another worktree, or concurrent Git access
	args.PrependOpcodes(
		&CheckoutUncached{Branch: existingPreviousBranch},
		&CheckoutUncached{Branch: cachedCurrentBranch},
	)
	return nil
}
