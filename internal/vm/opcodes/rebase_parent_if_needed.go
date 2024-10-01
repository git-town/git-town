package opcodes

import (
	"errors"
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type RebaseParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentIfNeeded) Run(args shared.RunArgs) error {
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
			var parentActiveInAnotherWorktree bool
			if parentBranchInfo, has := branchInfos.FindByLocalName(parent).Get(); has {
				parentActiveInAnotherWorktree = parentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
			} else {
				parentActiveInAnotherWorktree = false
			}
			program = append(program, &RebaseParent{
				CurrentBranch:               branch,
				ParentActiveInOtherWorktree: parentActiveInAnotherWorktree,
			})
			break
		}
		// here the parent isn't local --> sync with its tracking branch, then also sync the grandparent until we find a local ancestor
		parentTrackingName := parent.AtRemote(gitdomain.RemoteOrigin)
		// merge the parent's tracking branch
		program = append(program, &RebaseBranch{
			Branch: parentTrackingName.BranchName(),
		})
		// pull updates from the youngest local ancestor
		ancestors := args.Config.Config.Lineage.Ancestors(branch)
		slices.Reverse(ancestors) // youngest first now
		if youngestLocalAncestor, has := branchInfos.FirstLocal(ancestors).Get(); has {
			program = append(program, &RebaseBranch{
				Branch: youngestLocalAncestor.BranchName(),
			})
		}
		branch = parent
	}
	args.PrependOpcodes(program...)
	return nil
}
