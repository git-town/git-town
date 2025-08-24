package swap

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func swapUpdateProposalStackLineagesProgram(program Mutable[program.Program], proposalStackLineageArgs forge.ProposalStackLineageArgs) {
	// TODO: the proposal stack lineage builder executes its own GetProposalFn calls. Some of
	// these calls are duplicative because similar calls are completed upstream. A separate
	// enhancement is needed to reduce duplicative GetProposalFn calls through the forge
	// connectors.
	//
	// NOTE: NewProposalStackLineageBuilder method stores all the proposals found in the stack lineage
	// in map for fast retrieve through the GetProposals used below.
	tree, err := forge.NewProposalStackLineageTree(proposalStackLineageArgs)
	if err != nil {
		fmt.Printf("failed to update proposal stack lineage: %s\n", err.Error())
		return
	}
	for branch, proposal := range tree.BranchToProposal {
		program.Value.Add(&opcodes.ProposalUpdateLineage{
			Current:         branch,
			CurrentProposal: proposal,
			LineageTree:     MutableSome(tree),
		})
	}
}
