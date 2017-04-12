package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoAutomaticAbortOnError
	NoUndoStep
}

// CreateAbortStep returns the abort step for this step.
func (step FetchUpstreamStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step FetchUpstreamStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step FetchUpstreamStep) Run() error {
	return script.RunCommand("git", "fetch", "upstream")
}
