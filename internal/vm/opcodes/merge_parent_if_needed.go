package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// MergeParent merges the branch that at runtime is the parent branch of the given branch into the given branch.
type MergeParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentIfNeeded) Run(args shared.RunArgs) error {
	branch := self.Branch
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		return errors.New("BranchInfos not provided")
	}
	for {
		parent, hasParent := args.Config.Config.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if parentIsLocal {
			// parent is local --> sync the current branch with its local parent branch, then we are done
			parentActiveInAnotherWorktree := branchInfos.BranchIsActiveInAnotherWorktree(parent)
			var parentToMerge gitdomain.BranchName
			if parentActiveInAnotherWorktree {
				parentToMerge = parent.TrackingBranch().BranchName()
			} else {
				parentToMerge = parent.BranchName()
			}
			program = append(program, &MergeParent{
				Parent: parentToMerge,
			})
			break
		}
		// here the parent isn't local --> sync with its tracking branch, then also sync the grandparent until we find a local ancestor
		parentTrackingName := parent.AtRemote(gitdomain.RemoteOrigin)
		// merge the parent's tracking branch
		program = append(program, &MergeParent{
			Parent: parentTrackingName.BranchName(),
		})
		// continue climbing the ancestry chain until we find a local parent
		branch = parent
	}
	args.PrependOpcodes(program...)
	return nil
}