package opcodes

import (
	"errors"
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
	branch := self.CurrentBranch
	toAppend := []shared.Opcode{}
	for {
		parent, hasParent := args.Config.Config.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		if args.Git.BranchExists(args.Backend, parent) {
			// parent is local
			var parentActiveInAnotherWorktree bool
			branchInfos, has := args.InitialBranchesSnapshot.Get()
			if !has {
				return errors.New("initial branches snapshot not provided")
			}
			if parentBranchInfo, has := branchInfos.Branches.FindByLocalName(parent).Get(); has {
				parentActiveInAnotherWorktree = parentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
			} else {
				parentActiveInAnotherWorktree = false
			}
			toAppend = append(toAppend, &MergeParent{
				CurrentBranch:               branch,
				ParentActiveInOtherWorktree: parentActiveInAnotherWorktree,
			})
			break
		} else {
			// parent isn't local
			parentTrackingName := parent.AtRemote(gitdomain.RemoteOrigin)
			// merge the parent's tracking branch
			toAppend = append(toAppend, &MergeBranchNoEdit{
				Branch: parentTrackingName.BranchName(),
			})
			// pull updates from the youngest local ancestor
			ancestors := args.Config.Config.Lineage.Ancestors(branch)
			slices.Reverse(ancestors) // sort youngest first
			if youngestLocalAncestor, has := args.Git.FirstExistingBranch(args.Backend, ancestors...).Get(); has {
				toAppend = append(toAppend, &MergeBranchNoEdit{
					Branch: youngestLocalAncestor.BranchName(),
				})
			}
			branch = parent
		}
	}
	args.PrependOpcodes(toAppend...)
	return nil
}
