package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ProposalLineageCreate struct {
	Branch                  gitdomain.LocalBranchName
	ProposalLineageIn       configdomain.ProposalLineageIn
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalLineageCreate) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	findProposalFn, hasFindProposalFn := connector.FindProposalFn().Get()
	if !hasFindProposalFn {
		return fmt.Errorf("connector does not support finding proposals")
	}

	targetBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch)
	if targetBranch.IsNone() {
		return fmt.Errorf("current branch has no parent. Cannot find proposal to create lineage")
	}

	proposalData, err := findProposalFn(self.Branch, targetBranch.GetOrPanic())
	if err != nil {
		return err
	}

	if proposalData.IsNone() {
		return fmt.Errorf("current branch has no proposal")
	}

	builder := configdomain.NewProposalLineageBuilder(connector, args.Config.Value.MainAndPerennials()...)
	lineageInformation := args.Config.Value.NormalConfig.Lineage.BranchLineage(self.Branch)
	for _, curr := range lineageInformation {
		currParent := args.Config.Value.NormalConfig.Lineage.Parent(curr)
		builder.AddBranch(curr, currParent)
	}

	lineageAsString := builder.Build(self.Branch, self.ProposalLineageIn)

	switch self.ProposalLineageIn {
	case configdomain.ProposalLineageInTerminal:
		fmt.Print(lineageAsString)
		return nil
	case configdomain.ProposalLineageOperationInProposalBody:
		op := &ProposalUpdateBody{
			Proposal:    proposalData.GetOrPanic().Data,
			UpdatedBody: lineageAsString,
		}
		return op.Run(args)
	case configdomain.ProposalLineageOperationInProposalComment:
		// TODO: Implement soon
	default:
		return nil
	}
	return nil
}
