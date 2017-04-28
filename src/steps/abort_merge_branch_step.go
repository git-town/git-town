package steps

import "github.com/Originate/git-town/src/script"

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step AbortMergeBranchStep) Run() error {
	return script.RunCommand("git", "merge", "--abort")
}
