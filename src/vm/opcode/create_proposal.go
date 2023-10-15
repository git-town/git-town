package opcode

import (
	"github.com/git-town/git-town/v9/src/browser"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// CreateProposal creates a new pull request for the current branch.
type CreateProposal struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *CreateProposal) Run(args shared.RunArgs) error {
	parentBranch := args.Runner.Config.Lineage()[op.Branch]
	prURL, err := args.Connector.NewProposalURL(op.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Runner.Frontend.FrontendRunner, args.Runner.Backend)
	return nil
}
