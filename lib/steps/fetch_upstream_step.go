package steps

import (
	"github.com/Originate/git-town/lib/script"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoOpStep
}

// Run executes this step.
func (step FetchUpstreamStep) Run() error {
	return script.RunCommand("git", "fetch", "upstream")
}
