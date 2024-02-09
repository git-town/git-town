package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(list, branch.RemoteName, syncFeatureStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, parentOtherWorktree, syncFeatureStrategy)
}
