package steps

import (
	"github.com/Originate/git-town/lib/script"
)

type StashOpenChangesStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
}

func (step StashOpenChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step StashOpenChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step StashOpenChangesStep) CreateUndoStepBeforeRun() Step {
	return RestoreOpenChangesStep{}
}

func (step StashOpenChangesStep) Run() error {
	err := script.RunCommand("git", "add", "-A")
	if err != nil {
		return err
	}
	return script.RunCommand("git", "stash")
}
