package opcodes

import (
	"fmt"

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
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	proposalBodyUpdater, canUpdateProposalBody := connector.(forgedomain.ProposalBodyUpdater)
	if !canUpdateProposalBody {
		args.FinalMessages.Add(messages.ConnectorCannotUpdateProposalBody)
		return nil
	}
	fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB", self.Branch, args.Config.Value.NormalConfig.Lineage)
	proposals, err := proposalSearcher.SearchProposals(self.Branch)
	if err != nil {
		args.FinalMessages.Addf(messages.ProposalFindProblem, err.Error())
		return nil
	}
	fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", proposals)
	// TODO: here it tries to update the old proposals based on already updated lineage.
	// In this case, the lineage lists branch "beta" as a child of "main", but such a proposal doesn't exist.
	// Possible solutions:
	// 1. Update the embedded stack when updating the proposal target.
	//    This is brittle because at that time we are in the middle of applying changes.
	//    The resulting stack might reflect an intermediate state.
	// 2. Update all proposals that have beta as the source branch.
	//    This seems more robust, at the cost of doing possibly two separate updates.
	//    But these updates are different in nature (update target branch and update body),
	//    so it seems okay to do them separately.
	//    This also seems more correct, since all branches for "beta" need to get this update anyways.
	for _, proposal := range proposals {
		fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
		oldProposalBody := proposal.Data.Data().Body.GetOrZero()
		lineageSection := proposallineage.RenderSection(args.Config.Value.NormalConfig.Lineage, self.Branch, args.Config.Value.NormalConfig.Order, args.Connector)
		updatedProposalBody := proposallineage.UpdateProposalBody(oldProposalBody, lineageSection)
		fmt.Println("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", oldProposalBody, updatedProposalBody)
		if updatedProposalBody == oldProposalBody {
			continue
		}
		if err = proposalBodyUpdater.UpdateProposalBody(proposal.Data, updatedProposalBody); err != nil {
			args.FinalMessages.Addf(messages.ProposalBodyUpdateProblem, err.Error())
		}
	}
	return nil
}
