package opcodes

import (
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCreate creates a new proposal for the current branch.
type ProposalCreate struct {
	Branch        gitdomain.LocalBranchName
	MainBranch    gitdomain.LocalBranchName
	ProposalBody  Option[gitdomain.ProposalBody]
	ProposalTitle Option[gitdomain.ProposalTitle]
}

func (self *ProposalCreate) Run(args shared.RunArgs) error {
	parentBranch, hasParentBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		args.FinalMessages.Addf(messages.ProposalNoParent, self.Branch)
		return nil
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	if proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder); canFindProposals {
		existingProposalOpt, err := proposalFinder.FindProposal(self.Branch, parentBranch)
		if err != nil {
			args.FinalMessages.Addf(messages.ProposalFindProblem, err.Error())
			goto createProposal
		}
		if existingProposal, hasExistingProposal := existingProposalOpt.Get(); hasExistingProposal {
			args.PrependOpcodes(
				&BrowserOpen{
					URL: existingProposal.Data.Data().URL,
				},
			)
			return nil
		}
	}

createProposal:
	// TODO: create proposal with embedded lineage here. The lineage is loaded below.
	err := connector.CreateProposal(forgedomain.CreateProposalArgs{
		Branch:         self.Branch,
		FrontendRunner: args.Frontend,
		MainBranch:     self.MainBranch,
		ParentBranch:   parentBranch,
		ProposalBody:   self.ProposalBody,
		ProposalTitle:  self.ProposalTitle,
	})
	if err != nil {
		return err
	}

	if args.Config.Value.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		// TODO: remove this once we embed the lineage when creating the proposal
		args.PrependOpcodes(&ProposalUpdateLineage{
			Branch: self.Branch,
		})
	}
	return nil
}
