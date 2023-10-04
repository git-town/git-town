package step

import (
	"github.com/git-town/git-town/v9/src/browser"
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateProposal creates a new pull request for the current branch.
type CreateProposal struct {
	Branch domain.LocalBranchName
	Empty
}

func (step *CreateProposal) Run(args RunArgs) error {
	parentBranch := args.Runner.Config.Lineage()[step.Branch]
	prURL, err := args.Connector.NewProposalURL(step.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Runner.Frontend.FrontendRunner, args.Runner.Backend)
	return nil
}
