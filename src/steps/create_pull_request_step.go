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
	parentBranch := repo.GetParentBranch(step.BranchName)
	browsers.Open(driver.NewPullRequestURL(step.BranchName, parentBranch), repo.LoggingShell)
	return nil
}
