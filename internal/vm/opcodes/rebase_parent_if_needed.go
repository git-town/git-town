package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

type RebaseParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentIfNeeded) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	for branch := self.Branch; ; {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if parentIsLocal {
			var branchToRebase gitdomain.BranchName
			if branchInfos.BranchIsActiveInAnotherWorktree(parent) {
				branchToRebase = parent.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
			} else {
				branchToRebase = parent.BranchName()
			}
			program = append(program, &RebaseBranch{
				Branch: branchToRebase,
			})
			break
		}
		// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
		parentTrackingName := parent.AtRemote(args.Config.Value.NormalConfig.DevRemote)
		program = append(program, &RebaseBranch{
			Branch: parentTrackingName.BranchName(),
		})
		branch = parent
	}
	args.PrependOpcodes(program...)
	return nil
}
