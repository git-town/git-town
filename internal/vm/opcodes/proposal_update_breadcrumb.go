package opcodes

import (
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/messages"
	"github.com/git-town/git-town/v23/internal/proposallineage"
	"github.com/git-town/git-town/v23/internal/vm/shared"
)

type ProposalUpdateBreadcrumb struct {
	Branch gitdomain.LocalBranchName
}

func (self *ProposalUpdateBreadcrumb) Run(args shared.RunArgs) error {
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
	proposals, err := proposalSearcher.SearchProposals(self.Branch)
	if err != nil {
		args.FinalMessages.Addf(messages.ProposalFindProblem, err.Error())
		return nil
	}
	for _, proposal := range proposals {
		oldProposalBody := proposal.Data.Data().Body.GetOrZero()
		branchTypes := args.Config.Value.BranchesAndTypes(args.BranchInfos.LocalBranches().NamesLocalBranches())
		lineageSection := proposallineage.RenderSection(proposallineage.RenderSectionArgs{
			BranchTypes:   branchTypes,
			Breadcrumb:    args.Config.Value.NormalConfig.ProposalBreadcrumb,
			Connector:     args.Connector,
			CurrentBranch: self.Branch,
			Direction:     args.Config.Value.NormalConfig.ProposalBreadcrumbDirection,
			Excluded:      args.Config.Value.NormalConfig.ProposalBreadcrumbExcludeBranches,
			Lineage:       args.Config.Value.NormalConfig.Lineage,
			Order:         args.Config.Value.NormalConfig.Order,
		})
		updatedProposalBody := proposallineage.UpdateProposalBody(oldProposalBody, lineageSection)
		if updatedProposalBody == oldProposalBody {
			continue
		}
		if err = proposalBodyUpdater.UpdateProposalBody(proposal.Data, updatedProposalBody); err != nil {
			args.FinalMessages.Addf(messages.ProposalBodyUpdateProblem, err.Error())
		}
	}
	return nil
}
