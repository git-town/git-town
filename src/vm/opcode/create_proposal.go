package opcode

import (
	"github.com/git-town/git-town/v10/src/browser"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/vm/shared"
)

// CreateProposal creates a new pull request for the current branch.
type CreateProposal struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch := args.Runner.Config.Lineage(args.Runner.Config.RemoveLocalConfigValue)[self.Branch]
	prURL, err := args.Connector.NewProposalURL(self.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Runner.Frontend.FrontendRunner, args.Runner.Backend)
	return nil
}
