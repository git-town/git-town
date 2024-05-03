package opcodes

import (
	"github.com/git-town/git-town/v14/src/browser"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CreateProposal creates a new proposal for the current branch.
type CreateProposal struct {
	Branch     gitdomain.LocalBranchName
	MainBranch gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CreateProposal) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch := args.Runner.Config.Config.Lineage[self.Branch]
	prURL, err := args.Connector.NewProposalURL(self.Branch, parentBranch, self.MainBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Runner.Frontend.Runner, args.Runner.Backend.Runner)
	return nil
}
