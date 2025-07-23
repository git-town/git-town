package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ProposalCreateLineageProposalBody struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalCreateLineageProposalBody) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	_, hasFindProposal := connector.FindProposalFn().Get()
	if !hasFindProposal {
		return fmt.Errorf("connector does not support finding proposals")
	}

	builder := configdomain.NewProposalLineageBuilder(connector, args.Config.Value.MainAndPerennials()...)
	lineageInformation := args.Config.Value.NormalConfig.Lineage.BranchLineage(self.Branch)
	for _, curr := range lineageInformation {
		currParent := args.Config.Value.NormalConfig.Lineage.Parent(curr)
		builder.AddBranch(curr, currParent)
	}

	_ = builder.Build(self.Branch, configdomain.LineageDisplayLocationProposalBody)
	return nil
}
