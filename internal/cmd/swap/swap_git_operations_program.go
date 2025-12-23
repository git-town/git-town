package swap

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type swapGitOperationsProgramArgs struct {
	children    []swapBranch
	current     swapBranch
	grandParent gitdomain.LocalBranchName
	parent      swapBranch
	program     Mutable[program.Program]
}

func swapGitOperationsProgram(args swapGitOperationsProgramArgs) {
	// first update the current branch proposal (if there is one) target, so the proposal is not closed.
	if currentBranchProposal, currentBranchHasProposal := args.current.proposal.Get(); currentBranchHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.grandParent,
			OldBranch: args.parent.name,
			Proposal:  currentBranchProposal,
		})
	}
	args.program.Value.Add(
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.grandParent.BranchName(),
			CommitsToRemove:    args.parent.info.LocalBranchName().Location(),
		},
	)
	localBranch, HasLocalBranch := args.current.info.LocalName.Get()
	trackingBranch, hasTrackingBranch := args.current.info.RemoteName.Get()
	if HasLocalBranch && hasTrackingBranch {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: localBranch, ForceIfIncludes: true, TrackingBranch: trackingBranch},
		)
	}

	// next, update the parent proposal (if there is one), then rebase parent branch onto current
	if parentBranchProposal, parentBranchHasProposal := args.parent.proposal.Get(); parentBranchHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.current.name,
			OldBranch: args.grandParent,
			Proposal:  parentBranchProposal,
		})
	}
	args.program.Value.Add(
		&opcodes.Checkout{
			Branch: args.parent.info.LocalBranchName(),
		},
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: args.current.info.LocalBranchName().BranchName(),
			CommitsToRemove:    args.grandParent.Location(),
		},
	)
	if trackingBranch, hasTrackingBranch := args.parent.info.RemoteName.Get(); hasTrackingBranch {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{
				CurrentBranch:   args.parent.info.LocalBranchName(),
				ForceIfIncludes: true,
				TrackingBranch:  trackingBranch,
			},
		)
	}

	// Finally, update the child branches of current
	for _, child := range args.children {
		if childBranchProposal, childBranchHasProposal := child.proposal.Get(); childBranchHasProposal {
			args.program.Value.Add(&opcodes.ProposalUpdateTarget{
				NewBranch: args.parent.name,
				OldBranch: args.current.name,
				Proposal:  childBranchProposal,
			})
		}
		args.program.Value.Add(
			&opcodes.Checkout{
				Branch: child.name,
			},
		)
		oldBranchSHA := args.current.info.GetLocalOrRemoteSHA()
		args.program.Value.Add(
			&opcodes.RebaseOnto{
				BranchToRebaseOnto: args.parent.info.LocalBranchName().BranchName(),
				CommitsToRemove:    oldBranchSHA.Location(),
			},
		)
		if childTracking, childHasTracking := child.info.RemoteName.Get(); childHasTracking {
			args.program.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   child.name,
					ForceIfIncludes: true,
					TrackingBranch:  childTracking,
				},
			)
		}
	}
	args.program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.current.info.LocalBranchName()})
}
