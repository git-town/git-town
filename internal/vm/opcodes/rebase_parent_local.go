package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type RebaseParentLocal struct {
	Branch                  gitdomain.LocalBranchName
	ParentSHAPreviousRun    Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseParentLocal) Run(args shared.RunArgs) error {
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
			program = append(program, &RebaseBranch{
				Branch: branchToRebaseOnto,
			})
		}
	}
	return nil
}
