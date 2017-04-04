package steps

import (
	"log"

	"github.com/Originate/git-town/lib/script"
)

var fetched bool

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
	if !fetched {
		err := script.RunCommand("git", "fetch", "--prune")
		if err != nil {
			log.Fatal(err)
		}
		fetched = true
	}
	return nil
}
