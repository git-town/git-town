package sync

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type BranchProposalsProgramArgs struct {
	Program                  Mutable[program.Program]
	ProposalStackLineageArgs forge.ProposalStackLineageArgs
}

// BranchProposalsProgram syncs all given proposals.
func BranchProposalsProgram(branchesToSync configdomain.BranchesToSync, args BranchProposalsProgramArgs) {
	proposalStackLineageBuilder := forge.NewProposalStackLineageBuilder(args.ProposalStackLineageArgs)
	builder, hasBuilder := proposalStackLineageBuilder.Get()
	if !hasBuilder {
		return
	}

	for _, branch := range branchesToSync {
		proposal, hasProposal := builder.GetProposal(branch.BranchInfo.LocalBranchName()).Get()
		if !hasProposal {
			continue
		}
		args.Program.Value.Add(&opcodes.ProposalUpdateBody{
			Proposal: proposal,
			UpdatedBody: forge.ProposalBodyUpdateWithStackLineage(proposal.Data.Data().Body.GetOrDefault(), builder.Build(forge.ProposalStackLineageArgs{
				Connector:                args.ProposalStackLineageArgs.Connector,
				CurrentBranch:            branch.BranchInfo.LocalBranchName(),
				Lineage:                  args.ProposalStackLineageArgs.Lineage,
				MainAndPerennialBranches: args.ProposalStackLineageArgs.MainAndPerennialBranches,
			})),
		})
	}
}
