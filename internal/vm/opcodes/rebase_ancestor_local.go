package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// rebases a branch against a local ancestor branch
type RebaseAncestorLocal struct {
	Ancestor                gitdomain.LocalBranchName
	Branch                  gitdomain.LocalBranchName
	CommitsToRemove         Option[gitdomain.Location]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAncestorLocal) Run(args shared.RunArgs) error {
	fmt.Println("1111111111111111111111111111111111111111111 RebaseAncestorLocal")
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
	fmt.Println("1111111111111111111111111111111111 branchToRebaseOnto", branchToRebaseOnto)
	if commitsToRemove, hasCommitsToRemove := self.CommitsToRemove.Get(); hasCommitsToRemove {
		fmt.Println("111111111111111111111111111111111 commitsToRemove", commitsToRemove)
		args.PrependOpcodes(&RebaseOnto{
			BranchToRebaseOnto: branchToRebaseOnto,
			CommitsToRemove:    commitsToRemove,
		})
	} else {
		fmt.Println("111111111111111111111111111111111 no commitsToRemove", self.Branch, branchToRebaseOnto)
		isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, branchToRebaseOnto)
		if err != nil {
			return err
		}
		fmt.Println("1111111111111111111111111111111111 isInSync", isInSync)
		if !isInSync {
			args.PrependOpcodes(&RebaseBranch{
				Branch: branchToRebaseOnto,
			})
		}
	}
	return nil
}
