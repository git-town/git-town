package opcodes

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/vm/shared"
	. "github.com/git-town/git-town/v18/pkg/prelude"
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
			parentInfos, hasParentInfos := args.BranchInfos.Get()
			parentInfo, hasParentInfo := parentInfos.FindByLocalName(parent).Get()
			var opcode shared.Opcode
			if hasParentInfos && hasParentInfo && !branchInfos.BranchIsActiveInAnotherWorktree(parent) {
				opcode = &RebaseOntoKeepDeleted{
					BranchToRebaseOnto: parent,
					CommitsToRemove:    parentInfo.GetSHA().Location(),
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
