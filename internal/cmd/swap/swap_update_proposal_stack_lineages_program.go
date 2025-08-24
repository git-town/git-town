package swap

import (
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
	builder, hasBuilder := forge.NewProposalStackLineageBuilder(proposalStackLineageArgs).Get()
	if !hasBuilder {
		return
	}

	for _, proposal := range builder.GetProposals() {
		program.Value.Add(&opcodes.ProposalUpdateBody{
			Proposal: proposal,
			UpdatedBody: forge.ProposalBodyUpdateWithStackLineage(proposal.Data.Data().Body.GetOrDefault(), builder.Build(forge.ProposalStackLineageArgs{
				Connector:                proposalStackLineageArgs.Connector,
				CurrentBranch:            proposal.Data.Data().Source,
				Lineage:                  proposalStackLineageArgs.Lineage,
				MainAndPerennialBranches: proposalStackLineageArgs.MainAndPerennialBranches,
			})),
		})
	}
}
