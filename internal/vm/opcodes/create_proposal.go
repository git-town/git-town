package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/browser"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// CreateProposal creates a new proposal for the current branch.
type CreateProposal struct {
	Branch                  gitdomain.LocalBranchName
	MainBranch              gitdomain.LocalBranchName
	ProposalBody            gitdomain.ProposalBody
	ProposalTitle           gitdomain.ProposalTitle
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateProposal) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{self}
}

func (self *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch, hasParentBranch := args.Config.Config.Lineage.Parent(self.Branch).Get()
	if !hasParentBranch {
		return fmt.Errorf(messages.ProposalNoParent, self.Branch)
	}
	connector, hasConnector := args.Connector.Get()
	if !hasConnector {
		return hostingdomain.UnsupportedServiceError()
	}
	prURL, err := connector.NewProposalURL(self.Branch, parentBranch, self.MainBranch, self.ProposalTitle, self.ProposalBody)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Frontend, args.Backend)
	return nil
}
