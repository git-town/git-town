package steps

import (
	"log"

	"github.com/Originate/git-town/lib/script"
)

type FetchStep struct {
	NoAutomaticAbortOnError
}

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
	err := script.RunCommand("git", "fetch", "--prune")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
