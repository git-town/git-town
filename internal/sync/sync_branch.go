package sync

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchProgram syncs the given branch.
func BranchProgram(branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	localName, hasLocalName := branchInfo.LocalName.Get()
	if !hasLocalName {
		remoteName, hasRemoteName := branchInfo.RemoteName.Get()
		if !hasRemoteName {
			panic(fmt.Sprintf("branch %v has neither a local nor remote name", branchInfo))
		}
		localName = remoteName.LocalBranchName()
		args.Program.Value.Add(&opcodes.CheckoutIfNeeded{
			Branch: localName,
		})
	}
	switch {
	case branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote:
		deletedBranchProgram(args.Program, localName, args)
	case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// Git Town doesn't sync branches that are active in another worktree
	default:
		localBranchProgram(args.Program, branchInfo, firstCommitMessage, args)
	}
	args.Program.Value.Add(&opcodes.EndOfBranchProgram{})
}

type BranchProgramArgs struct {
	BranchInfos   gitdomain.BranchInfos
	Config        configdomain.ValidatedConfig
	InitialBranch gitdomain.LocalBranchName
	Program       Mutable[program.Program]
	PushBranches  configdomain.PushBranches
	Remotes       gitdomain.Remotes
}

// localBranchProgram provides the program to sync a local branch.
func localBranchProgram(list Mutable[program.Program], branch gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	localName, hasLocalName := branch.LocalName.Get()
	if !hasLocalName {
		return
	}
	isMainOrPerennialBranch := args.Config.IsMainOrPerennialBranch(localName)
	if isMainOrPerennialBranch && !args.Remotes.HasOrigin() {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	list.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localName})
	branchType := args.Config.BranchType(localName)
	fmt.Println("111111111111111111111111111111111111111", branchType)
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.Offline,
			program:            list,
			pushBranches:       args.PushBranches,
			remoteName:         branch.RemoteName,
			syncStrategy:       args.Config.SyncFeatureStrategy.SyncStrategy(),
		})
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(branch, args)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.InitialBranch, featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.Offline,
			program:            list,
			pushBranches:       args.PushBranches,
			remoteName:         branch.RemoteName,
			syncStrategy:       args.Config.SyncFeatureStrategy.SyncStrategy(),
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, branch)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(branch.RemoteName, args.Program)
	case configdomain.BranchTypePrototypeBranch:
		FeatureBranchProgram(featureBranchArgs{
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			offline:            args.Config.Offline,
			program:            list,
			pushBranches:       false,
			remoteName:         branch.RemoteName,
			syncStrategy:       args.Config.SyncPrototypeStrategy.SyncStrategy(),
		})
	}
	if args.PushBranches.IsTrue() && args.Remotes.HasOrigin() && args.Config.IsOnline() && branchType.ShouldPush(localName == args.InitialBranch) {
		switch {
		case !branch.HasTrackingBranch():
			list.Value.Add(&opcodes.CreateTrackingBranch{Branch: localName})
		case isMainOrPerennialBranch:
			list.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: localName})
		default:
			pushFeatureBranchProgram(list, localName, args.Config.SyncFeatureStrategy)
		}
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(args pullParentBranchOfCurrentFeatureBranchOpcodeArgs) {
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		args.program.Value.Add(&opcodes.MergeParentIfNeeded{
			Branch: args.branch,
		})
	case configdomain.SyncFeatureStrategyRebase:
		args.program.Value.Add(&opcodes.RebaseParentIfNeeded{
			Branch: args.branch,
		})
	case configdomain.SyncFeatureStrategyCompress:
		args.program.Value.Add(&opcodes.MergeParentIfNeeded{
			Branch: args.branch,
		})
	}
}

type pullParentBranchOfCurrentFeatureBranchOpcodeArgs struct {
	branch       gitdomain.LocalBranchName
	program      Mutable[program.Program]
	syncStrategy configdomain.SyncFeatureStrategy
}

func pushFeatureBranchProgram(list Mutable[program.Program], branch gitdomain.LocalBranchName, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	switch syncFeatureStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		list.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: branch})
	case configdomain.SyncFeatureStrategyRebase:
		list.Value.Add(&opcodes.ForcePushCurrentBranch{ForceIfIncludes: true})
	case configdomain.SyncFeatureStrategyCompress:
		list.Value.Add(&opcodes.ForcePushCurrentBranch{ForceIfIncludes: false})
	}
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(list Mutable[program.Program], otherBranch gitdomain.RemoteBranchName, strategy configdomain.SyncPerennialStrategy) {
	switch strategy {
	case configdomain.SyncPerennialStrategyMerge:
		list.Value.Add(&opcodes.Merge{Branch: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyRebase:
		list.Value.Add(&opcodes.RebaseBranch{Branch: otherBranch.BranchName()})
	}
}
