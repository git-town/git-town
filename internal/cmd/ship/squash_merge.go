package ship

import (
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

type shipDataMerge struct {
	authors []gitdomain.Author
	remotes gitdomain.Remotes
}

func determineMergeData(repo execute.OpenRepoResult, branch, parent gitdomain.LocalBranchName) (result shipDataMerge, err error) {
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return result, err
	}
	branchAuthors, err := repo.Git.BranchAuthors(repo.Backend, branch, parent)
	if err != nil {
		return result, err
	}
	return shipDataMerge{
		authors: branchAuthors,
		remotes: remotes,
	}, nil
}

func shipProgramSquashMerge(prog Mutable[program.Program], sharedData sharedShipData, squashMergeData shipDataMerge, commitMessage Option[gitdomain.CommitMessage]) {
	prog.Value.Add(&opcodes.BranchEnsureShippableChanges{Branch: sharedData.branchNameToShip, Parent: sharedData.targetBranchName})
	localTargetBranch, _ := sharedData.targetBranch.LocalName.Get()
	if sharedData.initialBranch != sharedData.targetBranchName {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	}
	if squashMergeData.remotes.HasDev(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.IsOnline() {
		UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	}
	prog.Value.Add(&opcodes.MergeSquashProgram{Authors: squashMergeData.authors, Branch: sharedData.branchNameToShip, CommitMessage: commitMessage, Parent: localTargetBranch})
	if squashMergeData.remotes.HasDev(sharedData.config.NormalConfig.DevRemote) && sharedData.config.NormalConfig.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchIfNeeded{CurrentBranch: sharedData.targetBranchName})
	}
	if !sharedData.dryRun {
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: sharedData.branchNameToShip})
	}
	if branchToShipRemoteName, hasRemoteName := sharedData.branchToShip.RemoteName.Get(); hasRemoteName {
		if sharedData.config.NormalConfig.IsOnline() {
			if sharedData.config.NormalConfig.ShipDeleteTrackingBranch {
				prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: branchToShipRemoteName})
			}
		}
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
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
}
