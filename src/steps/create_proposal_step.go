package steps

import (
	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreateProposalStep creates a new pull request for the current branch.
type CreateProposalStep struct {
	NoOpStep
	Branch string
}

func (step *CreateProposalStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	parentBranch := repo.Config.ParentBranch(step.Branch)
	prURL, err := driver.NewProposalURL(step.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, repo.LoggingShell)
	return nil
}
