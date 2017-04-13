package steps

import "github.com/Originate/git-town/lib/script"

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoExpectedError
	NoUndoStep
	Sha string
}

// Run executes this step.
func (step RevertCommitStep) Run() error {
	return script.RunCommand("git", "revert", step.Sha)
}
