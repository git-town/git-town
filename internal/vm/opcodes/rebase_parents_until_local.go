package opcodes

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

type RebaseParentsUntilLocal struct {
	Branch                  gitdomain.LocalBranchName
	PreviousSHA             Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentsUntilLocal) Run(args shared.RunArgs) error {
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
		parentIsPerennial := args.Config.Value.IsMainOrPerennialBranch(parent)
		if args.Detached.IsTrue() && parentIsPerennial {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if !parentIsLocal {
			// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
			parentTrackingName := parent.AtRemote(args.Config.Value.NormalConfig.DevRemote).BranchName()
			isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, parentTrackingName)
			if err != nil {
				return err
			}
			if !isInSync {
				program = append(program, &RebaseBranch{
					Branch: parentTrackingName,
				})
			}
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
		if previousParentSHA, hasPreviousParentSHA := self.PreviousSHA.Get(); hasPreviousParentSHA {
			program = append(program, &RebaseOntoKeepDeleted{
				BranchToRebaseOnto: branchToRebase,
				CommitsToRemove:    previousParentSHA.Location(),
				Upstream:           None[gitdomain.LocalBranchName](),
			})
		} else {
			isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, branchToRebase)
			if err != nil {
				return err
			}
			if !isInSync {
				program = append(program, &RebaseBranch{
					Branch: branchToRebase,
				})
			}
		}
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
