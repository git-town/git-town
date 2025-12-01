package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge"
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
		return fmt.Errorf(messages.ProposalNoParent, self.Branch)
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}

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
		if proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder); canFindProposals {
			lineageTree, err := forge.NewProposalStackLineageTree(forge.ProposalStackLineageArgs{
				Connector:                proposalFinder,
				CurrentBranch:            self.Branch,
				Lineage:                  args.Config.Value.NormalConfig.Lineage,
				MainAndPerennialBranches: args.Config.Value.MainAndPerennials(),
				Order:                    args.Config.Value.NormalConfig.Order,
			})
			if err != nil {
				// TODO: make sure error message return from failing to construct lineage is consistent across all invocations
				fmt.Printf("failed to construct proposal stack lineage: %s\n", err.Error())
			}
			proposalOpt, err := proposalFinder.FindProposal(self.Branch, parentBranch)
			if err != nil {
				return err
			}
			if proposalOpt.IsSome() {
				args.PrependOpcodes(&ProposalUpdateLineage{
					Current:         self.Branch,
					CurrentProposal: proposalOpt,
					LineageTree:     MutableSome(lineageTree),
				})
			}
		}
	}
	return nil
}
