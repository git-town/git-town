package opcodes

import (
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type ProposalUpdateLineage struct {
	Current         gitdomain.LocalBranchName
	CurrentProposal Option[forgedomain.Proposal]
	LineageTree     OptionalMutable[forge.ProposalStackLineageTree]
}

func (self *ProposalUpdateLineage) Run(args shared.RunArgs) error {
	proposal, hasProposal := self.CurrentProposal.Get()
	if !hasProposal {
		return nil
	}
	lineageArgs := forge.ProposalStackLineageArgs{
		Connector:                forgedomain.ProposalFinderFromConnector(args.Connector),
		CurrentBranch:            self.Current,
		Lineage:                  args.Config.Value.NormalConfig.Lineage,
		MainAndPerennialBranches: args.Config.Value.MainAndPerennials(),
		Order:                    args.Config.Value.NormalConfig.Order,
	}
	builder, hasBuilder := forge.NewProposalStackLineageBuilder(lineageArgs, self.LineageTree).Get()
	if !hasBuilder {
		return nil
	}
	if err := builder.UpdateStack(lineageArgs); err != nil {
		return err
	}
	args.PrependOpcodes(&ProposalUpdateBody{
		Proposal:    proposal,
		UpdatedBody: forge.ProposalBodyUpdateWithStackLineage(proposal.Data.Data().Body.GetOrZero(), builder.Build(lineageArgs)),
	})
	return nil
}
