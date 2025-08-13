package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases a branch against a local ancestor branch
type RebaseAncestorLocal struct {
	Ancestor                gitdomain.LocalBranchName
	Branch                  gitdomain.LocalBranchName
	CommitsToRemove         Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAncestorLocal) Run(args shared.RunArgs) error {
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	var branchToRebaseOnto gitdomain.BranchName
	if branchInfos.BranchIsActiveInAnotherWorktree(self.Ancestor) {
		branchToRebaseOnto = self.Ancestor.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
	} else {
		branchToRebaseOnto = self.Ancestor.BranchName()
	}
	if commitsToRemove, hasCommitsToRemove := self.CommitsToRemove.Get(); hasCommitsToRemove {
		// Here we rebase onto the new parent, while removing the commits that the parent had in the last run.
		// This removes old versions of commits that were amended by the user.
		// The new commits of the parent get added back during the rebase.
		args.PrependOpcodes(&RebaseOnto{
			BranchToRebaseOnto: branchToRebaseOnto,
			CommitsToRemove:    commitsToRemove.Location(),
		})
	} else {
		isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, branchToRebaseOnto)
		if err != nil {
			return err
		}
		if !isInSync {
			args.PrependOpcodes(&RebaseBranch{
				Branch: branchToRebaseOnto,
			})
		}
	}
	return nil
}
