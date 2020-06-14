package steps

import (
	"github.com/git-town/git-town/src/browsers"
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *CreatePullRequestStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	parentBranch := repo.Config.GetParentBranch(step.BranchName)
	prURL, err := driver.NewPullRequestURL(step.BranchName, parentBranch)
	if err != nil {
		return err
	}
	browsers.Open(prURL, repo.LoggingShell)
	return nil
}
