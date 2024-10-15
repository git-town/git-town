package ship

import (
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

func shipProgramFastForward(sharedData sharedShipData, squashMergeData shipDataMerge) program.Program {
	prog := NewMutable(&program.Program{})
	prog.Value.Add(&opcodes.EnsureHasShippableChanges{Branch: sharedData.branchNameToShip, Parent: sharedData.targetBranchName})
	if sharedData.initialBranch != sharedData.targetBranchName {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	}
	if squashMergeData.remotes.HasOrigin() && sharedData.config.Config.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	}
	prog.Value.Add(&opcodes.MergeFastForward{Branch: sharedData.branchNameToShip})
	if squashMergeData.remotes.HasOrigin() && sharedData.config.Config.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: sharedData.targetBranchName})
	}
	if !sharedData.dryRun {
		prog.Value.Add(&opcodes.BranchParentDelete{Branch: sharedData.branchNameToShip})
	}
	if branchToShipRemoteName, hasRemoteName := sharedData.branchToShip.RemoteName.Get(); hasRemoteName {
		if sharedData.config.Config.IsOnline() {
			if sharedData.config.Config.ShipDeleteTrackingBranch {
				prog.Value.Add(&opcodes.DeleteTrackingBranch{Branch: branchToShipRemoteName})
			}
		}
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageSetParentToGrandParent{Branch: child})
	}
	if !sharedData.isShippingInitialBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.initialBranch})
	}
	prog.Value.Add(&opcodes.BranchLocalDelete{Branch: sharedData.branchNameToShip})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   sharedData.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}
