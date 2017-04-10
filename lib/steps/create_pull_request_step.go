package steps

import (
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type CreatePullRequestStep struct {
	NoAutomaticAbort
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
	repository := git.GetUrlRepositoryName(git.GetRemoteOriginUrl())
	parentBranch := git.GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestUrl(repository, step.BranchName, parentBranch))
	return nil
}
