package sync

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
)

// BranchProgram syncs the given branch.
func BranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	parentBranchInfo := args.BranchInfos.FindByLocalName(args.Config.Lineage.Parent(branch.LocalName))
	parentOtherWorktree := parentBranchInfo != nil && parentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
	switch {
	case branch.SyncStatus == gitdomain.SyncStatusDeletedAtRemote:
		syncDeletedBranchProgram(args.Program, branch, parentOtherWorktree, args)
	case branch.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// Git Town doesn't sync branches that are active in another worktree
	default:
		ExistingBranchProgram(args.Program, branch, parentOtherWorktree, args)
	}
	args.Program.Add(&opcodes.EndOfBranchProgram{})
}

type BranchProgramArgs struct {
	BranchInfos gitdomain.BranchInfos
	Config      *configdomain.FullConfig
	Program     *program.Program
	PushBranch  bool
	Remotes     gitdomain.Remotes
}

// ExistingBranchProgram provides the opcode to sync a particular branch.
func ExistingBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	isFeatureBranch := args.Config.IsFeatureBranch(branch.LocalName)
	if !isFeatureBranch && !args.Remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&opcodes.Checkout{Branch: branch.LocalName})
	branchType := args.Config.BranchType(branch.LocalName)
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(list, branch, parentOtherWorktree, args.Config.SyncFeatureStrategy)
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branch, args)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branch, args)
	}
	if args.PushBranch && args.Remotes.HasOrigin() && args.Config.IsOnline() && branchType.ShouldPush() {
		switch {
		case !branch.HasTrackingBranch():
			list.Add(&opcodes.CreateTrackingBranch{Branch: branch.LocalName})
		case !isFeatureBranch:
			list.Add(&opcodes.PushCurrentBranch{CurrentBranch: branch.LocalName})
		default:
			pushFeatureBranchProgram(list, branch.LocalName, args.Config.SyncFeatureStrategy)
		}
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(list *program.Program, currentBranch gitdomain.LocalBranchName, parentOtherWorktree bool, strategy configdomain.SyncFeatureStrategy) {
	switch strategy {
	case configdomain.SyncFeatureStrategyMerge:
		list.Add(&opcodes.MergeParent{CurrentBranch: currentBranch, ParentActiveInOtherWorktree: parentOtherWorktree})
	case configdomain.SyncFeatureStrategyRebase:
		list.Add(&opcodes.RebaseParent{CurrentBranch: currentBranch, ParentActiveInOtherWorktree: parentOtherWorktree})
	}
}

// pullTrackingBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the remote branch of the current feature branch into the current feature branch.
func pullTrackingBranchOfCurrentFeatureBranchOpcode(list *program.Program, trackingBranch gitdomain.RemoteBranchName, strategy configdomain.SyncFeatureStrategy) {
	switch strategy {
	case configdomain.SyncFeatureStrategyMerge:
		list.Add(&opcodes.Merge{Branch: trackingBranch.BranchName()})
	case configdomain.SyncFeatureStrategyRebase:
		list.Add(&opcodes.RebaseBranch{Branch: trackingBranch.BranchName()})
	}
}

func pushFeatureBranchProgram(list *program.Program, branch gitdomain.LocalBranchName, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	switch syncFeatureStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		list.Add(&opcodes.PushCurrentBranch{CurrentBranch: branch})
	case configdomain.SyncFeatureStrategyRebase:
		list.Add(&opcodes.ForcePushCurrentBranch{})
	}
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(list *program.Program, otherBranch gitdomain.RemoteBranchName, strategy configdomain.SyncPerennialStrategy) {
	switch strategy {
	case configdomain.SyncPerennialStrategyMerge:
		list.Add(&opcodes.Merge{Branch: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyRebase:
		list.Add(&opcodes.RebaseBranch{Branch: otherBranch.BranchName()})
	}
}
