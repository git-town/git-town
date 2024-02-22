package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args featureBranchArgs) {
	if args.branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(args.program, args.branch.RemoteName, args.syncStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(args)
}

type featureBranchArgs struct {
	branch              gitdomain.BranchInfo             // the branch to sync
	parentOtherWorktree bool                             // whether the parent of this branch exists on another worktre
	program             *program.Program                 // the program to update
	syncStrategy        configdomain.SyncFeatureStrategy // the sync-feature-strategy
}

type SyncFeatureBranchArgs struct{}
