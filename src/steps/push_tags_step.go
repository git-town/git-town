package steps

import (
	"github.com/git-town/git-town/src/script"
)

// PushTagsStep pushes newly created Git tags to the remote.
type PushTagsStep struct {
	NoOpStep
}

// Run executes this step.
func (step *PushTagsStep) Run() error {
	return script.RunCommand("git", "push", "--tags")
}
