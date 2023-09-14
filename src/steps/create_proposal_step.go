package steps

import (
	"github.com/git-town/git-town/v9/src/browser"
	"github.com/git-town/git-town/v9/src/domain"
)

// CreateProposalStep creates a new pull request for the current branch.
type CreateProposalStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *CreateProposalStep) Run(args RunArgs) error {
	parentBranch := args.Run.Config.Lineage()[step.Branch]
	prURL, err := args.Connector.NewProposalURL(step.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, args.Run.Frontend.FrontendRunner, args.Run.Backend)
	return nil
}
