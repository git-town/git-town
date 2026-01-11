package ship

import (
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type shipDataMerge struct {
	authors []gitdomain.Author
	remotes gitdomain.Remotes
}

func determineMergeData(repo execute.OpenRepoResult, branch, parent gitdomain.LocalBranchName) (shipDataMerge, error) {
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return shipDataMerge{}, err
	}
	branchAuthors, err := repo.Git.BranchAuthors(repo.Backend, branch, parent)
	return shipDataMerge{
		authors: branchAuthors,
		remotes: remotes,
	}, err
}

func shipProgramSquashMerge(prog Mutable[program.Program], repo execute.OpenRepoResult, sharedData sharedShipData, squashMergeData shipDataMerge, commitMessage Option[gitdomain.CommitMessage]) {
	prog.Value.Add(&opcodes.BranchEnsureShippableChanges{Branch: sharedData.branchToShip, Parent: sharedData.targetBranchName})
	localTargetBranch, _ := sharedData.targetBranch.LocalName().Get()
	if sharedData.initialBranch != sharedData.targetBranchName {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	}
	if squashMergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.Offline.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	}
	prog.Value.Add(&opcodes.MergeSquashProgram{Authors: squashMergeData.authors, Branch: sharedData.branchToShip, CommitMessage: commitMessage, Parent: localTargetBranch})
	if squashMergeData.remotes.HasRemote(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.Offline.IsOnline() {
		if trackingBranch, hasTrackingBranch := sharedData.targetBranch.RemoteName.Get(); hasTrackingBranch {
			prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: sharedData.targetBranchName, TrackingBranch: trackingBranch})
		}
	}
	if branchToShipRemoteName, hasRemoteName := sharedData.branchToShipInfo.RemoteName.Get(); hasRemoteName {
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
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: sharedData.branchToShip})
	}
	if !sharedData.isShippingInitialBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.initialBranch})
	}
	prog.Value.Add(&opcodes.BranchLocalDelete{Branch: sharedData.branchToShip})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   repo.UnvalidatedConfig.NormalConfig.DryRun,
		InitialStashSize:         sharedData.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
}
