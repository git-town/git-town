package opcodes

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
	fmt.Println("111111111111111111111111111111111111111111111111111111111111111111111111111")
	spew.Dump(connector)
	proposalSearcher, canSearchProposals := connector.(forgedomain.ProposalSearcher)
	if !canSearchProposals {
		fmt.Println("444444444444444444444444444444444444444444444444444444444444444444444444444", proposalSearcher)
		args.FinalMessages.Add(messages.ConnectorCannotSearchProposals)
		return nil
	}
	fmt.Println("333333333333333333333333333333333333333333333333333333333333333333333333333")
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
		fmt.Println("222222222222222222222222222222222222222222222222222222222222222222222222222", proposal)
		oldProposalBody := proposal.Data.Data().Body.GetOrZero()
		lineageSection := proposallineage.RenderSection(args.Config.Value.NormalConfig.Lineage, self.Branch, args.Config.Value.NormalConfig.Order, args.Connector)
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
