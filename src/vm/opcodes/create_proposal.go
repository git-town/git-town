package opcodes

import (
	"github.com/git-town/git-town/v14/src/browser"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// CreateProposal creates a new proposal for the current branch.
type CreateProposal struct {
	Branch                  gitdomain.LocalBranchName
	MainBranch              gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CreateProposal) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		self,
	}
}

func (self *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch := args.Config.Config.Lineage[self.Branch]
	if args.Connector == nil {
		return hostingdomain.UnsupportedServiceError()
	}
	prURL, err := args.Connector.NewProposalURL(self.Branch, parentBranch, self.MainBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Frontend.Runner, args.Backend.Runner)
	return nil
}
