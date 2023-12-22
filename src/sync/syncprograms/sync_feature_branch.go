package syncprograms

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// SyncFeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func SyncFeatureBranchProgram(list *program.Program, branch syncdomain.BranchInfo, parentOtherWorktree bool, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(list, branch.RemoteName, syncFeatureStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, parentOtherWorktree, syncFeatureStrategy)
}
