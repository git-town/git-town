package steps

import "github.com/Originate/git-town/lib/script"

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
	Sha string
}

// CreateAbortStep returns the abort step for this step.
func (step RevertCommitStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step RevertCommitStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step RevertCommitStep) Run() error {
	return script.RunCommand("git", "revert", step.Sha)
}
