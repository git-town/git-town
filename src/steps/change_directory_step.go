package steps

import (
	"log"
	"os"

	"github.com/Originate/git-town/src/script"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	NoOpStep
	Directory string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step ChangeDirectoryStep) CreateUndoStepBeforeRun() Step {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return ChangeDirectoryStep{Directory: dir}
}

// Run executes this step.
func (step ChangeDirectoryStep) Run() error {
	_, err := os.Stat(step.Directory)
	if err == nil {
		script.PrintCommand("cd", step.Directory)
		return os.Chdir(step.Directory)
	}
	return nil
}
