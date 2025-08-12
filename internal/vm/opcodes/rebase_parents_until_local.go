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
	fmt.Println("22222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222")
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
		fmt.Println("3333333333333333333333333333333333333333333333333333333333333333333333", parent)
		parentIsPerennial := args.Config.Value.IsMainOrPerennialBranch(parent)
		if args.Detached.IsTrue() && parentIsPerennial {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if !parentIsLocal {
			fmt.Println("4444444444444444444444444444444444444444444444444444444444444444444444444444444444444444444 PARENT IS NOT LOCAL")
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
		if parentSHAPreviousRun, hasParentSHAPreviousRun := self.ParentSHAPreviousRun.Get(); hasParentSHAPreviousRun {
			// Here we rebase onto the new parent, while removing the commits that the parent had in the last run.
			// This removes old versions of commits that were amended by the user.
			// The new commits of the parent get added back during the rebase.
			fmt.Println("55555555555555555555555555555555555555555555555555555555555 HAS PREVIOUS RUN", branchToRebase)
			program = append(program, &RebaseOnto{
				BranchToRebaseOnto: branchToRebaseOnto,
				CommitsToRemove:    parentSHAPreviousRun.Location(),
			})
		} else {
			isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, branchToRebaseOnto)
			if err != nil {
				return err
			}
			if !isInSync {
				fmt.Println("666666666666666666666666666666666666666666666666 NOT IN SYNC")
				program = append(program, &RebaseBranch{
					Branch: branchToRebaseOnto,
				})
			}
		}
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
