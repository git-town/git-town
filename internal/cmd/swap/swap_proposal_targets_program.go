package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type swapProposalTargetsProgramArg struct {
	branchToSwap      swapBranch
	grandParentBranch gitdomain.LocalBranchName
	parentBranch      swapBranch
	program           Mutable[program.Program]
}

func swapProposalTargetsProgram(args swapProposalTargetsProgramArg) {
	branchToSwapProposal, branchToSwapHasProposal := args.branchToSwap.proposal.Get()
	parentBranchProposal, parentBranchHasProposal := args.parentBranch.proposal.Get()

	if !branchToSwapHasProposal && !parentBranchHasProposal {
		return
	}

	if branchToSwapHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.grandParentBranch,
			OldBranch: args.parentBranch.name,
			Proposal:  branchToSwapProposal,
		})
	}

	if parentBranchHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.branchToSwap.name,
			OldBranch: args.grandParentBranch,
			Proposal:  parentBranchProposal,
		})
	}
}
