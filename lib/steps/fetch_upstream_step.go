package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type FetchUpstreamStep struct{}

func (step FetchUpstreamStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step FetchUpstreamStep) Run() error {
	return script.RunCommand("git", "fetch", "upstream")
}

func (step FetchUpstreamStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
