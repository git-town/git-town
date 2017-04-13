package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// RestoreOpenChangesStep restores stashed away changes into the workspace.
type RestoreOpenChangesStep struct {
	NoOpStep
}

// Run executes this step.
func (step RestoreOpenChangesStep) Run() error {
	return script.RunCommand("git", "stash", "pop")
}
