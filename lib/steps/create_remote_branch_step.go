package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type CreateRemoteBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	BranchName string
	Sha        string
}

func (step CreateRemoteBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CreateRemoteBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CreateRemoteBranchStep) Run() error {
	return script.RunCommand("git", "push", "origin", step.Sha+":refs/heads/"+step.BranchName)
}
