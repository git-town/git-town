package steps

import (
	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// CreatePullRequestStep creates a new pull request for the current branch.
type CreatePullRequestStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step CreatePullRequestStep) Run() error {
	driver := drivers.GetCodeHostingDriver()
	repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
	parentBranch := git.GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestURL(repository, step.BranchName, parentBranch))
	return nil
}
