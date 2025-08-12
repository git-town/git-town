package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases a branch against a local ancestor branch
type RebaseAncestorLocal struct {
	Ancestor                gitdomain.LocalBranchName
	Branch                  gitdomain.LocalBranchName
	CommitsToRemove         Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

// TODO: new sync architecture
//
// Currently, Git Town syncs a branch like this:
//  1. rebase the branch onto main to remove all old commits.
//     This is unnecessarily reductive, we should rebase onto its new parent and remove old commits.
//  2. rebase the branch against its new parent
//
// It should be possible to perform both in one operation: rebase onto the new parent, removing the old commits in the branch
// The old commits are determined this way:
//   - original parent deleted during this sync --> SHA of the local parent branch
//   - original parent not deleted --> SHA of that branch at the last run
func (self *RebaseAncestorLocal) Run(args shared.RunArgs) error {
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	var branchToRebaseOnto gitdomain.BranchName
	if branchInfos.BranchIsActiveInAnotherWorktree(self.Ancestor) {
		branchToRebaseOnto = self.Ancestor.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
	} else {
		branchToRebaseOnto = self.Ancestor.BranchName()
	}
	if commitsToRemove, hasCommitsToRemove := self.CommitsToRemove.Get(); hasCommitsToRemove {
		// Here we rebase onto the new parent, while removing the commits that the parent had in the last run.
		// This removes old versions of commits that were amended by the user.
		// The new commits of the parent get added back during the rebase.
		args.PrependOpcodes(&RebaseOnto{
			BranchToRebaseOnto: branchToRebaseOnto,
			CommitsToRemove:    commitsToRemove.Location(),
		})
	} else {
		isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, branchToRebaseOnto)
		if err != nil {
			return err
		}
		if !isInSync {
			args.PrependOpcodes(&RebaseBranch{
				Branch: branchToRebaseOnto,
			})
		}
	}
	return nil
}
