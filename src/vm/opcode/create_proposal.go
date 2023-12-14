package opcode

import (
	"github.com/git-town/git-town/v11/src/browser"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// CreateProposal creates a new proposal for the current branch.
type CreateProposal struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CreateProposal) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch := args.Runner.GitTown.Lineage(args.Runner.GitTown.RemoveLocalConfigValue)[self.Branch]
	prURL, err := args.Connector.NewProposalURL(self.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Runner.Frontend.FrontendRunner, args.Runner.Backend)
	return nil
}
