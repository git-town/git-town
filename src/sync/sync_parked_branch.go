package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// PerennialBranchProgram adds the opcodes to sync the observed branch with the given name.
func ParkedBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	if branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(list, branch.RemoteName, syncFeatureStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(list, branch.LocalName, parentOtherWorktree, syncFeatureStrategy)
}
