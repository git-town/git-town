package ship

import (
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func shipProgramFastForward(prog Mutable[program.Program], repo execute.OpenRepoResult, sharedData sharedShipData, squashMergeData shipDataMerge) {
	prog.Value.Add(&opcodes.BranchEnsureShippableChanges{Branch: sharedData.branchNameToShip, Parent: sharedData.targetBranchName})
	if sharedData.initialBranch != sharedData.targetBranchName {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	}
	if squashMergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.Offline.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	}
	prog.Value.Add(&opcodes.MergeFastForward{Branch: sharedData.branchNameToShip.BranchName()})
	if squashMergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.Offline.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: sharedData.targetBranchName})
	}
	if branchToShipRemoteName, hasRemoteName := sharedData.branchToShip.RemoteName.Get(); hasRemoteName {
		if sharedData.config.NormalConfig.Offline.IsOnline() {
			if sharedData.config.NormalConfig.ShipDeleteTrackingBranch {
				prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: branchToShipRemoteName})
			}
		}
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	if !repo.UnvalidatedConfig.NormalConfig.DryRun {
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: sharedData.branchNameToShip})
	}
	if !sharedData.isShippingInitialBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.initialBranch})
	}
	prog.Value.Add(&opcodes.BranchLocalDelete{Branch: sharedData.branchNameToShip})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   repo.UnvalidatedConfig.NormalConfig.DryRun,
		InitialStashSize:         sharedData.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
}
