package opcodes

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/proposallineage2"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalUpdateLineage struct {
	Branch   gitdomain.LocalBranchName
	Proposal Option[forgedomain.Proposal]
}

func (self *ProposalUpdateLineage) Run(args shared.RunArgs) error {
	proposal, hasProposal := self.Proposal.Get()
	if !hasProposal {
		return nil
	}
	lineageSection := proposallineage2.RenderSection(args.Config.Value.NormalConfig.Lineage, self.Branch, args.Config.Value.NormalConfig.Order, args.Connector)
	lineageArgs := proposallineage.ProposalStackLineageArgs{
		Connector:                forgedomain.ProposalFinderFromConnector(args.Connector),
		CurrentBranch:            self.Branch,
		Lineage:                  args.Config.Value.NormalConfig.Lineage,
		MainAndPerennialBranches: args.Config.Value.MainAndPerennials(),
		Order:                    args.Config.Value.NormalConfig.Order,
	}
	builder, hasBuilder := proposallineage.NewBuilder(lineageArgs, self.LineageTree).Get()
	if !hasBuilder {
		return nil
	}
	args.PrependOpcodes(&ProposalUpdateBody{
		Proposal:    proposal,
		UpdatedBody: proposallineage.Add(proposal.Data.Data().Body.GetOrZero(), builder.Build(lineageArgs)),
	})
	return nil
}
