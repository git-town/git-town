package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
	currentBranchProposal, currentBranchHasProposal := args.current.proposal.Get()
	if currentBranchHasProposal {
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
	if args.current.info.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.current.info.LocalBranchName(), ForceIfIncludes: true},
		)
	}

	// next, update the parent proposal (if there is one), then rebase parent branch onto current
	parentBranchProposal, parentBranchHasProposal := args.parent.proposal.Get()
	if parentBranchHasProposal {
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
	if args.parent.info.HasTrackingBranch() {
		args.program.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: args.parent.info.LocalBranchName(), ForceIfIncludes: true},
		)
	}

	// Finally, update the child branches of current
	for _, child := range args.children {
		childBranchProposal, childBranchHasProposal := child.proposal.Get()
		if childBranchHasProposal {
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
		oldBranchSHA, hasOldBranchSHA := args.current.info.LocalSHA.Get()
		if !hasOldBranchSHA {
			oldBranchSHA = args.current.info.RemoteSHA.GetOrDefault()
		}
		args.program.Value.Add(
			&opcodes.RebaseOnto{
				BranchToRebaseOnto: args.parent.info.LocalBranchName().BranchName(),
				CommitsToRemove:    oldBranchSHA.Location(),
			},
		)
		if child.info.HasTrackingBranch() {
			args.program.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   child.name,
					ForceIfIncludes: true,
				},
			)
		}
	}
	args.program.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.current.info.LocalBranchName()})
}
