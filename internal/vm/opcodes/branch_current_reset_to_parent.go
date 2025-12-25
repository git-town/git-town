package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchCurrentResetToParent resets all commits in the current branch to the parent.
type BranchCurrentResetToParent struct {
	CurrentBranch gitdomain.LocalBranchName
}

func (self *BranchCurrentResetToParent) Run(args shared.RunArgs) error {
	parentName, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.CurrentBranch).Get()
	if !hasParent {
		return nil
	}
	parentInfo, hasParentInfo := args.BranchInfos.FindLocalOrRemote(parentName).Get()
	if !hasParentInfo {
		return fmt.Errorf(messages.BranchInfoNotFound, parentName)
	}
	args.PrependOpcodes(&BranchReset{Target: parentInfo.GetLocalOrRemoteName()})
	return nil
}
