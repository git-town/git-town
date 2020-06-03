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
	Driver     drivers.CodeHostingDriver
}

// Run executes this step.
func (step *CreatePullRequestStep) Run(repo *git.ProdRepo) error {
	parentBranch := repo.GetParentBranch(step.BranchName)
	script.OpenBrowser(step.Driver.GetNewPullRequestURL(step.BranchName, parentBranch))
	return nil
}
