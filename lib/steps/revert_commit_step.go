package steps

import "github.com/Originate/git-town/lib/script"

type RevertCommitStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	Sha string
}

func (step RevertCommitStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) Run() error {
	return script.RunCommand("git", "revert", step.Sha)
}
