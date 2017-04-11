package steps

import "github.com/Originate/git-town/lib/script"

type CreateBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	BranchName    string
	StartingPoint string
}

func (step CreateBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CreateBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
