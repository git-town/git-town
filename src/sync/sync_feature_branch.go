package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args SyncFeatureBranchArgs) {
	if args.branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(args.list, args.branch.RemoteName, args.syncFeatureStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(args.list, args.branch.LocalName, args.parentOtherWorktree, args.syncFeatureStrategy)
}

type SyncFeatureBranchArgs struct {
	list                *program.Program
	branch              gitdomain.BranchInfo
	parentOtherWorktree bool
	syncFeatureStrategy configdomain.SyncFeatureStrategy
}
