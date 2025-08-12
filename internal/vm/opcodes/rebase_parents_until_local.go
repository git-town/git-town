package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type RebaseParentsUntilLocal struct {
	Branch                  gitdomain.LocalBranchName
	ParentSHAPreviousRun    Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentsUntilLocal) Run(args shared.RunArgs) error {
	fmt.Println("RebaseParentsUntilLocal")
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
		var branchToRebaseOnto gitdomain.BranchName
		if branchInfos.BranchIsActiveInAnotherWorktree(parent) {
			branchToRebaseOnto = parent.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
		} else {
			branchToRebaseOnto = parent.BranchName()
		}
		program = append(program, &RebaseOnto{
			BranchToRebaseOnto: branchToRebaseOnto,
			CommitsToRemove:    commitsToRemove,
			Upstream:           None[gitdomain.LocalBranchName](),
		})
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
