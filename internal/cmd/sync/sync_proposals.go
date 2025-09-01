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
	builder, hasBuilder := forge.NewProposalStackLineageBuilder(args.ProposalStackLineageArgs).Get()
	if !hasBuilder {
		return
	}

	for _, branch := range branchesToSync {
		// TODO: there are now multiple places that load and use proposals for branches.
		// To avoid double-loading the same proposal data in one run,
		// extract an object that caches the already known proposals,
		// i.e. which branch has which proposal,
		// and loads missing proposal info on demand.
		proposal, hasProposal := builder.GetProposal(branch.BranchInfo.LocalBranchName()).Get()
		if !hasProposal {
			continue
		}
		args.Program.Value.Add(&opcodes.ProposalUpdateBody{
			Proposal: proposal,
			UpdatedBody: forge.ProposalBodyUpdateWithStackLineage(proposal.Data.Data().Body.GetOrZero(), builder.Build(forge.ProposalStackLineageArgs{
				Connector:                args.ProposalStackLineageArgs.Connector,
				CurrentBranch:            branch.BranchInfo.LocalBranchName(),
				Lineage:                  args.ProposalStackLineageArgs.Lineage,
				MainAndPerennialBranches: args.ProposalStackLineageArgs.MainAndPerennialBranches,
			})),
		})
	}
}
