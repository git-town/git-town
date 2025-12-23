package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ProposalCreate creates a new proposal for the current branch.
type ProposalCreateIfNeeded struct {
	Branch        gitdomain.LocalBranchName
	MainBranch    gitdomain.LocalBranchName
	ProposalBody  Option[gitdomain.ProposalBody]
	ProposalTitle Option[gitdomain.ProposalTitle]
}

func (self *ProposalCreateIfNeeded) Run(args shared.RunArgs) error {
	parentBranch, hasParentBranch := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		args.FinalMessages.Add(fmt.Sprintf(messages.ProposalNoParent, self.Branch))
		return nil
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return forgedomain.UnsupportedServiceError()
	}
	proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder)
	if !canFindProposals {
		return nil
	}
	existingProposal, err := proposalFinder.FindProposal(self.Branch, parentBranch)
	if err != nil {
		return err
	}
	if existingProposal.IsSome() {
		return nil
	}
	args.PrependOpcodes(&ProposalCreate{
		Branch:        self.Branch,
		MainBranch:    self.MainBranch,
		ProposalBody:  self.ProposalBody,
		ProposalTitle: self.ProposalTitle,
	})
	return nil
}
