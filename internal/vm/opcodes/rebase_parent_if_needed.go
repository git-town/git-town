package opcodes

import (
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/vm/shared"
)

type RebaseParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	PreviousSHA             Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentIfNeeded) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	branch := self.Branch
	for {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if parentIsLocal {
			var opcode shared.Opcode
			previousBranchInfos, hasPreviousBranchInfos := args.PreviousBranchInfos.Get()
			previousParentInfo, hasPreviousParentInfo := previousBranchInfos.FindByLocalName(parent).Get()
			if hasPreviousBranchInfos && hasPreviousParentInfo && !branchInfos.BranchIsActiveInAnotherWorktree(parent) {
				opcode = &RebaseOntoKeepDeleted{
					BranchToRebaseOnto: parent,
					CommitsToRemove:    previousParentInfo.GetLocalOrRemoteSHA().Location(),
					Upstream:           None[gitdomain.LocalBranchName](),
				}
			} else {
				opcode = &RebaseBranch{
					Branch: parent.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName(),
				}
			}
			program = append(program, opcode)
			break
		}
		program = append(program, &RebaseBranch{
			Branch: branchToRebase,
		})
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
