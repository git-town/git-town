package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type FetchUpstreamStep struct {
	NoAutomaticAbort
}

func (step FetchUpstreamStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step FetchUpstreamStep) Run() error {
	return script.RunCommand("git", "fetch", "upstream")
}
