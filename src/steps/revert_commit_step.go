package steps

import "github.com/Originate/git-town/src/script"

// RevertCommitStep reverts the commit with the given sha.
type RevertCommitStep struct {
	NoOpStep
	Sha string
}

// Run executes this step.
func (step RevertCommitStep) Run() error {
	return script.RunCommand("git", "revert", step.Sha)
}
