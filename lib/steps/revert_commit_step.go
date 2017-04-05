package steps

import "github.com/Originate/git-town/lib/script"

type RevertCommitStep struct{}

func (step RevertCommitStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step RevertCommitStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step RevertCommitStep) Run() error {
	return script.RunCommand("git", "revert", "HEAD")
}

func (step RevertCommitStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
