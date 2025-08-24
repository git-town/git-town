package opcodes

import (
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type ProposalUpdateLineage struct {
	Current                 gitdomain.LocalBranchName
	CurrentProposal         Option[forgedomain.Proposal]
	LineageTree             OptionalMutable[forge.ProposalStackLineageTree]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateLineage) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return nil
	}

	proposal, hasProposal := self.CurrentProposal.Get()
	if !hasProposal {
		return nil
	}

	lineageArgs := forge.ProposalStackLineageArgs{
		Connector:                connector,
		CurrentBranch:            self.Current,
		Lineage:                  args.Config.Value.NormalConfig.Lineage,
		MainAndPerennialBranches: args.Config.Value.MainAndPerennials(),
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
		UpdatedBody: forge.ProposalBodyUpdateWithStackLineage(proposal.Data.Data().Body.GetOrDefault(), builder.Build(lineageArgs)),
	})
	return nil
}
