package sync

import (
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

// BranchProgram syncs the given branch.
func BranchProgram(localName gitdomain.LocalBranchName, branchInfo gitdomain.BranchInfo, firstCommitMessage Option[gitdomain.CommitMessage], args BranchProgramArgs) {
	parentNameOpt := args.Config.NormalConfig.Lineage.Parent(localName)
	parentName, hasParentName := parentNameOpt.Get()
	parentSHAInitial := None[gitdomain.SHA]()
	if hasParentName {
		if parentBranchInfo, hasParentBranchInfo := args.BranchInfos.FindLocalOrRemote(parentName, args.Config.NormalConfig.DevRemote).Get(); hasParentBranchInfo {
			parentSHAInitial = parentBranchInfo.LocalSHA.Or(parentBranchInfo.RemoteSHA)
		}
	}
	usesRebaseSyncStrategy := args.Config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase
	ancestorToRemove, hasAncestorToRemove := args.Config.NormalConfig.Lineage.YoungestAncestorWithin(localName, args.BranchesToDelete.Value.Values()).Get()
	parentSHAPrevious := None[gitdomain.SHA]()
	if parent, has := parentNameOpt.Get(); has {
		if branchInfosLastRun, has := args.BranchInfosPrevious.Get(); has {
			if parentInfoLastRun, has := branchInfosLastRun.FindByLocalName(parent).Get(); has {
				parentSHAPrevious = Some(parentInfoLastRun.GetLocalOrRemoteSHA())
			}
		}
	}
	trackingBranchGone := branchInfo.SyncStatus == gitdomain.SyncStatusDeletedAtRemote
	hasDescendents := args.Config.NormalConfig.Lineage.HasDescendents(localName)
	switch {
	case hasAncestorToRemove && ancestorToRemove == parentName && trackingBranchGone && hasDescendents:
		args.BranchesToDelete.Value.Add(localName)
	case hasAncestorToRemove && ancestorToRemove == parentName:
		if usesRebaseSyncStrategy {
			RemoveAncestorCommits(RemoveAncestorCommitsArgs{
				Ancestor:          ancestorToRemove.BranchName(),
				Branch:            localName,
				HasTrackingBranch: branchInfo.HasTrackingBranch(),
				Program:           args.Program,
				RebaseOnto:        args.Config.ValidatedConfigData.MainBranch, // TODO: RebaseOnto the latest existing parent, which isn't always main
			})
		}
	case usesRebaseSyncStrategy && trackingBranchGone && hasDescendents:
		args.BranchesToDelete.Value.Add(localName)
	case trackingBranchGone:
		deletedBranchProgram(localName, parentNameOpt, parentSHAInitial, parentSHAPrevious, args)
	case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
		// cannot sync branches that are active in another worktree
	default:
		if hasAncestorToRemove && usesRebaseSyncStrategy {
			RemoveAncestorCommits(RemoveAncestorCommitsArgs{
				Ancestor:          ancestorToRemove.BranchName(),
				Branch:            localName,
				HasTrackingBranch: branchInfo.HasTrackingBranch(),
				Program:           args.Program,
				RebaseOnto:        parentName,
			})
		}
		localBranchProgram(localBranchProgramArgs{
			BranchProgramArgs:  args,
			branchInfo:         branchInfo,
			firstCommitMessage: firstCommitMessage,
			localName:          localName,
			parentNameInitial:  parentNameOpt,
			parentSHAInitial:   parentSHAInitial,
			parentSHAPrevious:  parentSHAPrevious,
		})
	}
	args.Program.Value.Add(&opcodes.ProgramEndOfBranch{})
}

type BranchProgramArgs struct {
	BranchInfos         gitdomain.BranchInfos                       // the initial BranchInfos, after "git fetch" ran
	BranchInfosPrevious Option[gitdomain.BranchInfos]               // the BranchInfos at the end of the previous Git Town command
	BranchesToDelete    Mutable[set.Set[gitdomain.LocalBranchName]] // branches that should be deleted after the branches are all synced
	Config              config.ValidatedConfig
	InitialBranch       gitdomain.LocalBranchName
	PrefetchBranchInfos gitdomain.BranchInfos // BranchInfos before "git fetch" ran
	Program             Mutable[program.Program]
	Prune               configdomain.Prune
	PushBranches        configdomain.PushBranches
	Remotes             gitdomain.Remotes
}

type localBranchProgramArgs struct {
	BranchProgramArgs
	branchInfo         gitdomain.BranchInfo
	firstCommitMessage Option[gitdomain.CommitMessage]
	localName          gitdomain.LocalBranchName
	parentNameInitial  Option[gitdomain.LocalBranchName]
	parentSHAInitial   Option[gitdomain.SHA]
	parentSHAPrevious  Option[gitdomain.SHA]
}

// localBranchProgram provides the program to sync a local branch.
func localBranchProgram(args localBranchProgramArgs) {
	branchType := args.Config.BranchType(args.localName)
	isMainOrPerennialBranch := branchType == configdomain.BranchTypeMainBranch || branchType == configdomain.BranchTypePerennialBranch
	if isMainOrPerennialBranch && !args.Remotes.HasRemote(args.Config.NormalConfig.DevRemote) {
		// perennial branch but no remote --> this branch cannot be synced
		return
	}
	args.Program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.localName})
	switch branchType {
	case configdomain.BranchTypeFeatureBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage:   args.firstCommitMessage,
			initialParentName:    args.parentNameInitial,
			initialParentSHA:     args.parentSHAInitial,
			localName:            args.localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: args.parentSHAPrevious,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         args.PushBranches,
			trackingBranch:       args.branchInfo.RemoteName,
		})
	case configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch:
		PerennialBranchProgram(args.branchInfo, args.BranchProgramArgs)
	case configdomain.BranchTypeParkedBranch:
		ParkedBranchProgram(args.Config.NormalConfig.SyncFeatureStrategy.SyncStrategy(), args.InitialBranch, featureBranchArgs{
			firstCommitMessage:   args.firstCommitMessage,
			initialParentName:    args.parentNameInitial,
			initialParentSHA:     args.parentSHAInitial,
			localName:            args.localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: args.parentSHAPrevious,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         args.PushBranches,
			trackingBranch:       args.branchInfo.RemoteName,
		})
	case configdomain.BranchTypeContributionBranch:
		ContributionBranchProgram(args.Program, args.branchInfo)
	case configdomain.BranchTypeObservedBranch:
		ObservedBranchProgram(args.branchInfo, args.Program)
	case configdomain.BranchTypePrototypeBranch:
		FeatureBranchProgram(args.Config.NormalConfig.SyncPrototypeStrategy.SyncStrategy(), featureBranchArgs{
			firstCommitMessage:   args.firstCommitMessage,
			initialParentName:    args.parentNameInitial,
			initialParentSHA:     args.parentSHAInitial,
			localName:            args.localName,
			offline:              args.Config.NormalConfig.Offline,
			parentSHAPreviousRun: args.parentSHAPrevious,
			program:              args.Program,
			prune:                args.Prune,
			pushBranches:         configdomain.PushBranches(args.branchInfo.HasTrackingBranch()),
			trackingBranch:       args.branchInfo.RemoteName,
		})
	}
	if args.PushBranches.ShouldPush() && args.Remotes.HasRemote(args.Config.NormalConfig.DevRemote) && args.Config.NormalConfig.Offline.IsOnline() && branchType.ShouldPush(args.localName == args.InitialBranch) {
		isMainBranch := branchType == configdomain.BranchTypeMainBranch
		switch {
		case !args.branchInfo.HasTrackingBranch():
			args.Program.Value.Add(&opcodes.BranchTrackingCreateIfLocalExists{Branch: args.localName})
		case isMainBranch && args.Remotes.HasUpstream() && args.Config.NormalConfig.SyncUpstream.ShouldSyncUpstream():
			if trackingBranch, hasTrackingBranch := args.branchInfo.RemoteName.Get(); hasTrackingBranch {
				args.Program.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: args.localName, TrackingBranch: trackingBranch})
			}
		case isMainOrPerennialBranch && !shouldPushPerennialBranch(args.branchInfo.SyncStatus):
			// don't push if its a perennial branch that doesn't need pushing
		case isMainOrPerennialBranch:
			if trackingBranch, hasTrackingBranch := args.branchInfo.RemoteName.Get(); hasTrackingBranch {
				args.Program.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: args.localName, TrackingBranch: trackingBranch})
			}
		default:
			if trackingBranch, hasTrackingBranch := args.branchInfo.RemoteName.Get(); hasTrackingBranch {
				pushFeatureBranchProgram(args.Program, args.localName, trackingBranch, args.Config.NormalConfig.SyncFeatureStrategy)
			}
		}
	}
}

// pullParentBranchOfCurrentFeatureBranchOpcode adds the opcode to pull updates from the parent branch of the current feature branch into the current feature branch.
func pullParentBranchOfCurrentFeatureBranchOpcode(args pullParentBranchOfCurrentFeatureBranchOpcodeArgs) {
	switch args.syncStrategy {
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.program.Value.Add(&opcodes.SyncFeatureBranchMerge{
			Branch:            args.branch,
			InitialParentName: args.parentNameInitial,
			InitialParentSHA:  args.parentSHAInitial,
			TrackingBranch:    args.trackingBranch,
		})
	case configdomain.SyncFeatureStrategyRebase:
		args.program.Value.Add(&opcodes.RebaseAncestorsUntilLocal{
			Branch:          args.branch,
			CommitsToRemove: args.parentSHAPrevious,
		})
	}
}

type pullParentBranchOfCurrentFeatureBranchOpcodeArgs struct {
	branch            gitdomain.LocalBranchName
	parentNameInitial Option[gitdomain.LocalBranchName]
	parentSHAInitial  Option[gitdomain.SHA]
	parentSHAPrevious Option[gitdomain.SHA]
	program           Mutable[program.Program]
	syncStrategy      configdomain.SyncFeatureStrategy
	trackingBranch    Option[gitdomain.RemoteBranchName]
}

func pushFeatureBranchProgram(prog Mutable[program.Program], branch gitdomain.LocalBranchName, trackingBranch gitdomain.RemoteBranchName, syncFeatureStrategy configdomain.SyncFeatureStrategy) {
	switch syncFeatureStrategy {
	case configdomain.SyncFeatureStrategyMerge:
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: branch})
	case configdomain.SyncFeatureStrategyRebase:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: branch, ForceIfIncludes: true, TrackingBranch: trackingBranch})
	case configdomain.SyncFeatureStrategyCompress:
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: branch, ForceIfIncludes: false, TrackingBranch: trackingBranch})
	}
}

func RemoveAncestorCommits(args RemoveAncestorCommitsArgs) {
	args.Program.Value.Add(
		&opcodes.CheckoutIfNeeded{Branch: args.Branch},
	)
	if args.HasTrackingBranch {
		args.Program.Value.Add(
			&opcodes.PullCurrentBranch{},
		)
	}
	args.Program.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.RebaseOnto.BranchName(),
			CommitsToRemove:    args.Ancestor.Location(),
		},
	)
	if args.HasTrackingBranch {
		args.Program.Value.Add(
			&opcodes.PushCurrentBranchForce{ForceIfIncludes: false},
		)
	}
}

type RemoveAncestorCommitsArgs struct {
	Ancestor          gitdomain.BranchName
	Branch            gitdomain.LocalBranchName
	HasTrackingBranch bool
	Program           Mutable[program.Program]
	RebaseOnto        gitdomain.LocalBranchName
}

func shouldPushPerennialBranch(syncStatus gitdomain.SyncStatus) bool {
	switch syncStatus {
	case
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind,
		gitdomain.SyncStatusLocalOnly,
		gitdomain.SyncStatusNotInSync:
		return true
	case
		gitdomain.SyncStatusDeletedAtRemote,
		gitdomain.SyncStatusOtherWorktree,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusUpToDate:
	}
	return false
}

// updateCurrentPerennialBranchOpcode provides the opcode to update the current perennial branch with changes from the given other branch.
func updateCurrentPerennialBranchOpcode(prog Mutable[program.Program], otherBranch gitdomain.RemoteBranchName, strategy configdomain.SyncPerennialStrategy) {
	switch strategy {
	case configdomain.SyncPerennialStrategyMerge:
		prog.Value.Add(&opcodes.MergeIntoCurrentBranch{BranchToMerge: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyRebase:
		prog.Value.Add(&opcodes.RebaseBranch{Branch: otherBranch.BranchName()})
	case configdomain.SyncPerennialStrategyFFOnly:
		prog.Value.Add(&opcodes.MergeFastForward{Branch: otherBranch.BranchName()})
	}
}
