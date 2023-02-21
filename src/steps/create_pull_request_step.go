package steps

import (
	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	Branch string
}

func (step *CreatePullRequestStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	parentBranch := repo.Config.ParentBranch(step.Branch)
	prURL, err := connector.NewChangeRequestURL(step.Branch, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, repo.LoggingShell)
	return nil
}
