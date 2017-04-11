package steps

import (
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type CreatePullRequestStep struct {
	BranchName string
}

func (step CreatePullRequestStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CreatePullRequestStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CreatePullRequestStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step CreatePullRequestStep) Run() error {
	driver := drivers.GetCodeHostingDriver()
	repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
	parentBranch := git.GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestURL(repository, step.BranchName, parentBranch))
	return nil
}
