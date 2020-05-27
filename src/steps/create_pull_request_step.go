package steps

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *CreatePullRequestStep) Run(repo *git.ProdRepo) error {
	driver := drivers.GetActiveDriver()
	parentBranch := git.Config().GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestURL(step.BranchName, parentBranch))
	return nil
}
