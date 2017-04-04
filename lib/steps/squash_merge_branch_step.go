package steps

import "github.com/Originate/git-town/lib/script"

type SquashMergeBranchStep struct {
	BranchName string
}

func (step SquashMergeBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step SquashMergeBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step SquashMergeBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step SquashMergeBranchStep) Run() error {
	return script.RunCommand("git", "merge", "--squash", step.BranchName)
}

func (step SquashMergeBranchStep) ShouldAbortOnError() (bool, string) {
	return false, ""
}
