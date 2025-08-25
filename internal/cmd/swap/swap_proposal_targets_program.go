package swap

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type swapProposalTargetsProgramArg struct {
	current     swapBranch
	grandParent gitdomain.LocalBranchName
	parent      swapBranch
	program     Mutable[program.Program]
}

func swapProposalTargetsProgram(args swapProposalTargetsProgramArg) {
	currentBranchProposal, currentBranchHasProposal := args.current.proposal.Get()
	parentBranchProposal, parentBranchHasProposal := args.parent.proposal.Get()
	if !currentBranchHasProposal && !parentBranchHasProposal {
		return
	}
	if currentBranchHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.grandParent,
			OldBranch: args.parent.name,
			Proposal:  currentBranchProposal,
		})
	}
	if parentBranchHasProposal {
		args.program.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: args.current.name,
			OldBranch: args.grandParent,
			Proposal:  parentBranchProposal,
		})
	}
}
