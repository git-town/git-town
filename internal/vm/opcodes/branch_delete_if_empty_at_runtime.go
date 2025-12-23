package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchDeleteIfEmptyAtRuntime deletes the given branch if it has no content at runtime.
type BranchDeleteIfEmptyAtRuntime struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchDeleteIfEmptyAtRuntime) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return nil
	}
	hasUnmergedChanges, err := args.Git.BranchHasUnmergedChanges(args.Backend, self.Branch, parent)
	if err != nil {
		return err
	}
	if hasUnmergedChanges {
		return nil
	}
	opcodes := []shared.Opcode{
		&CheckoutDescendentOrOther{
			Branch: self.Branch,
		},
	}
	if branchInfo, hasBranchInfo := args.BranchInfos.FindByLocalName(self.Branch).Get(); hasBranchInfo {
		if trackingBranch, hasTrackingBranch := branchInfo.RemoteName.Get(); hasTrackingBranch {
			opcodes = append(opcodes, &BranchTrackingDelete{Branch: trackingBranch})
		}
	}
	opcodes = append(opcodes,
		&BranchLocalDeleteContent{
			BranchToDelete:     self.Branch,
			BranchToRebaseOnto: args.Config.Value.ValidatedConfigData.MainBranch,
		},
		&LineageBranchRemove{
			Branch: self.Branch,
		},
	)
	args.PrependOpcodes(opcodes...)
	return nil
}
