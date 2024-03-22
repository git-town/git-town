package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// FeatureBranchProgram adds the opcodes to sync the feature branch with the given name.
func FeatureBranchProgram(args featureBranchArgs) {
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		syncFeatureBranchMergeProgram(args)
	case configdomain.SyncFeatureStrategyRebase:
		syncFeatureBranchRebaseProgram(args)
	}
}

type featureBranchArgs struct {
	branch              gitdomain.BranchInfo             // the branch to sync
	offline             configdomain.Offline             // whether offline mode is enabled
	parentOtherWorktree bool                             // whether the parent of this branch exists on another worktre
	program             *program.Program                 // the program to update
	syncStrategy        configdomain.SyncFeatureStrategy // the sync-feature-strategy
}

// syncs the given feature branch using the "merge" sync strategy
func syncFeatureBranchMergeProgram(args featureBranchArgs) {
	if args.branch.HasTrackingBranch() {
		pullTrackingBranchOfCurrentFeatureBranchOpcode(args.program, args.branch.RemoteName, args.syncStrategy)
	}
	pullParentBranchOfCurrentFeatureBranchOpcode(args)
}

// syncs the given feature branch using the "rebase" sync strategy
func syncFeatureBranchRebaseProgram(args featureBranchArgs) {
	// rebase against parent
	args.program.Add(&opcodes.RebaseParent{
		CurrentBranch:               args.branch.LocalName,
		ParentActiveInOtherWorktree: args.parentOtherWorktree,
	})
	if args.branch.HasTrackingBranch() && !args.offline.Bool() {
		args.program.Add(&opcodes.RebaseFeatureTrackingBranch{RemoteBranch: args.branch.RemoteName})
	}
}
