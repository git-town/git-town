package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// RebaseAncestorLocal rebases a branch against a local ancestor branch.
type RebaseAncestorLocal struct {
	Ancestor        gitdomain.LocalBranchName
	Branch          gitdomain.LocalBranchName
	CommitsToRemove Option[gitdomain.SHA]
}

func (self *RebaseAncestorLocal) Run(args shared.RunArgs) error {
	branchToRebaseOnto := self.Ancestor.BranchName()
	commitsToRemove, hasCommitsToRemove := self.CommitsToRemove.Get()
	ancestorSHA := None[gitdomain.SHA]()
	if hasCommitsToRemove {
		sha, err := args.Git.SHAForBranch(args.Backend, branchToRebaseOnto)
		if err != nil {
			return err
		}
		ancestorSHA = Some(sha)
	}
	if hasCommitsToRemove && !self.CommitsToRemove.Equal(ancestorSHA) {
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
