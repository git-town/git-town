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
		if !parentIsLocal {
			// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
			parentTrackingName := parent.AtRemote(args.Config.Value.NormalConfig.DevRemote)
			program = append(program, &RebaseBranch{
				Branch: parentTrackingName.BranchName(),
			})
			branch = parent
			continue
		}
		// here the parent is local
		var branchToRebase gitdomain.BranchName
		if branchInfos.BranchIsActiveInAnotherWorktree(parent) {
			branchToRebase = parent.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
		} else {
			branchToRebase = parent.BranchName()
		}
		var opcode shared.Opcode
		if previousParentSHA, hasPreviousParentSHA := self.PreviousSHA.Get(); hasPreviousParentSHA {
			opcode = &RebaseOntoKeepDeleted{
				BranchToRebaseOnto: branchToRebase,
				CommitsToRemove:    previousParentSHA.Location(),
				Upstream:           None[gitdomain.LocalBranchName](),
			}
		} else {
			opcode = &RebaseBranch{
				Branch: branchToRebase,
			}
		}
		program = append(program, opcode)
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
