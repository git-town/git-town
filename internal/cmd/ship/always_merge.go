package ship

import (
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type shipProgramAlwaysMergeArgs struct {
	commitMessage Option[gitdomain.CommitMessage]
	mergeData     shipDataMerge
	prog          Mutable[program.Program]
	sharedData    sharedShipData
}

func shipProgramAlwaysMerge(repo execute.OpenRepoResult, args shipProgramAlwaysMergeArgs) {
	args.prog.Value.Add(&opcodes.BranchEnsureShippableChanges{Branch: args.sharedData.branchToShip, Parent: args.sharedData.targetBranchName})
	if args.sharedData.initialBranch != args.sharedData.targetBranchName {
		args.prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.sharedData.targetBranchName})
	}
	if args.mergeData.remotes.HasRemote(args.sharedData.config.NormalConfig.DevRemote) && args.sharedData.config.NormalConfig.Offline.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(args.prog.Value, args.sharedData.proposalsOfChildBranches)
	}
	args.prog.Value.Add(&opcodes.MergeAlwaysProgram{Branch: args.sharedData.branchToShip, CommitMessage: args.commitMessage})
	if args.mergeData.remotes.HasRemote(args.sharedData.config.NormalConfig.DevRemote) && args.sharedData.config.NormalConfig.Offline.IsOnline() {
		if trackingBranch, hasTrackingBranch := args.sharedData.branchToShipInfo.RemoteName.Get(); hasTrackingBranch {
			args.prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: args.sharedData.targetBranchName, TrackingBranch: trackingBranch})
		}
	}
	if branchToShipRemoteName, hasRemoteName := args.sharedData.branchToShipInfo.RemoteName.Get(); hasRemoteName {
		if args.sharedData.config.NormalConfig.Offline.IsOnline() {
			if args.sharedData.config.NormalConfig.ShipDeleteTrackingBranch {
				args.prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: branchToShipRemoteName})
			}
		}
	}
	for _, child := range args.sharedData.childBranches {
		args.prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	args.prog.Value.Add(&opcodes.LineageParentRemove{Branch: args.sharedData.branchToShip})
	if !args.sharedData.isShippingInitialBranch {
		args.prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.sharedData.initialBranch})
	}
	args.prog.Value.Add(&opcodes.BranchLocalDelete{Branch: args.sharedData.branchToShip})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{args.sharedData.previousBranch}
	cmdhelpers.Wrap(args.prog, cmdhelpers.WrapOptions{
		DryRun:                   repo.UnvalidatedConfig.NormalConfig.DryRun,
		InitialStashSize:         args.sharedData.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         !args.sharedData.isShippingInitialBranch && args.sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
}
