package opcodes

import (
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParentIfNeeded struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentIfNeeded) Run(args shared.RunArgs) error {
	if parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get(); hasParent {
		if args.Git.BranchExists(args.Backend, parent) {
			// parent is local
			var parentActiveInAnotherWorktree bool
			if parentBranchInfo, has := args.InitialBranchesSnapshot.Branches.FindByLocalName(parent).Get(); has {
				parentActiveInAnotherWorktree = parentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
			} else {
				parentActiveInAnotherWorktree = false
			}
			args.PrependOpcodes(&MergeParent{
				CurrentBranch:               self.CurrentBranch,
				ParentActiveInOtherWorktree: parentActiveInAnotherWorktree,
			})
		} else {
			// parent isn't local
			parentTrackingName := parent.AtRemote(gitdomain.RemoteOrigin)
			// pull updates from the youngest local ancestor
			ancestors := args.Config.Config.Lineage.Ancestors(self.CurrentBranch)
			slices.Reverse(ancestors) // sort youngest first
			if youngestLocalAncestor, has := args.Git.FirstExistingBranch(args.Backend, ancestors...).Get(); has {
				args.PrependOpcodes(&MergeBranchNoEdit{
					Branch: youngestLocalAncestor.BranchName(),
				})
			}
			// merge the parent's tracking branch
			args.PrependOpcodes(&MergeBranchNoEdit{
				Branch: parentTrackingName.BranchName(),
			})
		}
	}
	return nil
}
