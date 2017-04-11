package steps

import "github.com/Originate/git-town/lib/script"

type FetchStep struct{}

func (step FetchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step FetchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step FetchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step FetchStep) Run() error {
	return script.RunCommand("git", "fetch", "--prune")
}
