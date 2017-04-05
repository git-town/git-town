package steps

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/drivers"
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

func (step CreatePullRequestStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step CreatePullRequestStep) Run() error {
	driver := drivers.GetCodeHostingDriver()
	repository := config.GetUrlRepositoryName(config.GetRemoteOriginUrl())
	parentBranch := config.GetParentBranch(step.BranchName)
	script.OpenBrowser(driver.GetNewPullRequestUrl(repository, step.BranchName, parentBranch))
	return nil
}

func (step CreatePullRequestStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
