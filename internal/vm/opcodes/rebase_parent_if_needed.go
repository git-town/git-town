package opcodes

import (
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/vm/shared"
	. "github.com/git-town/git-town/v19/pkg/prelude"
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
			if previousBranchInfos, has := args.PreviousBranchInfos.Get(); has {
				if previousParentInfo, has := previousBranchInfos.FindByLocalName(parent).Get(); has {
					if !branchInfos.BranchIsActiveInAnotherWorktree(parent) {
						opcode = &RebaseOntoKeepDeleted{
							BranchToRebaseOnto: parent,
							CommitsToRemove:    previousParentInfo.GetLocalOrRemoteSHA().Location(),
							Upstream:           None[gitdomain.LocalBranchName](),
						}
					}
				}
			}
			if opcode == nil {
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
