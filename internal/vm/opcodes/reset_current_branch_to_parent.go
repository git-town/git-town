package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ResetCurrentBranch resets all commits in the current branch.
type ResetCurrentBranchToParent struct {
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ResetCurrentBranchToParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Config.Lineage.Parent(self.CurrentBranch).Get()
	if !hasParent {
		return nil
	}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		return errors.New(messages.BranchInfosNotProvided)
	}
	parentIsLocal := branchInfos.HasLocalBranch(parent)
	var target gitdomain.BranchName
	if parentIsLocal {
		target = parent.BranchName()
	} else {
		target = parent.TrackingBranch().BranchName()
	}
	args.PrependOpcodes(&ResetBranch{Target: target})
	return nil
}
