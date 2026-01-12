package opcodes

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ProposalUpdateLineage struct {
	Branch gitdomain.LocalBranchName
}

func (self *ProposalUpdateLineage) Run(args shared.RunArgs) error {
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		args.FinalMessages.Add(forgedomain.UnsupportedServiceError().Error())
		return nil
	}
	proposalSearcher, canSearchProposals := connector.(forgedomain.ProposalSearcher)
	if !canSearchProposals {
		args.FinalMessages.Add(messages.ConnectorCannotSearchProposals)
		return nil
	}
	proposalBodyUpdater, canUpdateProposalBody := connector.(forgedomain.ProposalBodyUpdater)
	if !canUpdateProposalBody {
		args.FinalMessages.Add(messages.ConnectorCannotUpdateProposalBody)
		return nil
	}
	parentBranch, hasParentBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		return nil
	}
	proposalOpt, err := proposalFinder.FindProposal(self.Branch, parentBranch)
	if err != nil {
		args.FinalMessages.Addf(messages.ProposalFindProblem, err.Error())
		return nil
	}
	proposal, hasProposal := proposalOpt.Get()
	if !hasProposal {
		return nil
	}
	oldProposalBody := proposal.Data.Data().Body.GetOrZero()
	lineageSection := proposallineage.RenderSection(args.Config.Value.NormalConfig.Lineage, self.Branch, args.Config.Value.NormalConfig.Order, args.Connector)
	updatedProposalBody := proposallineage.UpdateProposalBody(oldProposalBody, lineageSection)
	if updatedProposalBody == oldProposalBody {
		return nil
	}
	if err = proposalBodyUpdater.UpdateProposalBody(proposal.Data, updatedProposalBody); err != nil {
		args.FinalMessages.Addf(messages.ProposalBodyUpdateProblem, err.Error())
	}
	return nil
}
