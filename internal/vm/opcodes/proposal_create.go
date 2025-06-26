package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// ProposalCreate creates a new proposal for the current branch.
type ProposalCreate struct {
	Branch                  gitdomain.LocalBranchName
	MainBranch              gitdomain.LocalBranchName
	ProposalBody            Option[gitdomain.ProposalBody]
	ProposalTitle           Option[gitdomain.ProposalTitle]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
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
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAaaa", connector)
	return connector.CreateProposal(forgedomain.CreateProposalArgs{
		Branch:         self.Branch,
		FrontendRunner: args.Frontend,
		MainBranch:     self.MainBranch,
		ParentBranch:   parentBranch,
		ProposalBody:   self.ProposalBody,
		ProposalTitle:  self.ProposalTitle,
	})
}
