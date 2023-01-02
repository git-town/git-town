package steps

import (
	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

func (step *CreatePullRequestStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	parentBranch := repo.Config.Ancestry.Parent(step.BranchName)
	prURL, err := driver.NewPullRequestURL(step.BranchName, parentBranch)
	if err != nil {
		return err
	}
	browser.Open(prURL, repo.LoggingShell)
	return nil
}
