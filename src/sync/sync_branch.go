package sync

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
)

// BranchProgram syncs the given branch.
func BranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	parentOtherWorktree := false
	if localName, hasLocalName := branch.LocalName.Get(); hasLocalName {
		parent, hasParent := args.Config.Lineage.Parent(localName).Get()
		if hasParent {
			parentBranchInfo, hasParentBranchInfo := args.BranchInfos.FindByLocalName(parent).Get()
			parentOtherWorktree = hasParentBranchInfo && parentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		}
	}
	switch {
	case branch.SyncStatus == gitdomain.SyncStatusDeletedAtRemote:
		syncDeletedBranchProgram(args.Program, branch.LocalName, parentOtherWorktree, args)
	case branch.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// Git Town doesn't sync branches that are active in another worktree
	default:
		ExistingBranchProgram(args.Program, branch, parentOtherWorktree, args)
	}
	args.Program.Add(&opcodes.EndOfBranchProgram{})
}

type BranchProgramArgs struct {
	BranchInfos   gitdomain.BranchInfos
	Config        configdomain.ValidatedConfig
	InitialBranch gitdomain.LocalBranchName
	Program       *program.Program
	PushBranch    bool
	Remotes       gitdomain.Remotes
}

// ExistingBranchProgram provides the opcode to sync a particular branch.
func ExistingBranchProgram(list *program.Program, branch gitdomain.BranchInfo, parentOtherWorktree bool, args BranchProgramArgs) {
	localName, hasLocalName := branch.LocalName.Get()
	if !hasLocalName {
		return
	}
	isMainOrPerennialBranch := args.Config.IsMainOrPerennialBranch(localName)
	if isMainOrPerennialBranch && !args.Remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Add(&opcodes.Checkout{Branch: localName})
	branchType := args.Config.BranchType(localName)
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(featureBranchArgs{
			branch:              branch,
			offline:             args.Config.Offline,
			parentOtherWorktree: parentOtherWorktree,
			program:             list,
			syncStrategy:        args.Config.SyncFeatureStrategy,
		})
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branch, args)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.InitialBranch, featureBranchArgs{
			branch:              branch,
			offline:             args.Config.Offline,
			parentOtherWorktree: parentOtherWorktree,
			program:             list,
			syncStrategy:        args.Config.SyncFeatureStrategy,
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, branch)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branch, args.Program)
	}
	if args.PushBranch && args.Remotes.HasOrigin() && args.Config.IsOnline() && branchType.ShouldPush(localName, args.InitialBranch) {
		switch {
		case !branch.HasTrackingBranch():
			list.Add(&opcodes.CreateTrackingBranch{Branch: localName})
		case isMainOrPerennialBranch:
			list.Add(&opcodes.PushCurrentBranch{CurrentBranch: localName})
		default:
			pushFeatureBranchProgram(list, localName, args.Config.SyncFeatureStrategy)
		}
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(args pullParentBranchOfCurrentFeatureBranchOpcodeArgs) {
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		args.program.Add(&opcodes.MergeParent{CurrentBranch: args.branch, ParentActiveInOtherWorktree: args.parentOtherWorktree})
	case configdomain.SyncFeatureStrategyRebase:
		args.program.Add(&opcodes.RebaseParent{CurrentBranch: args.branch, ParentActiveInOtherWorktree: args.parentOtherWorktree})
	}
}

type pullParentBranchOfCurrentFeatureBranchOpcodeArgs struct {
	branch              gitdomain.LocalBranchName
	program             *program.Program
	parentOtherWorktree bool
	syncStrategy        configdomain.SyncFeatureStrategy
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
