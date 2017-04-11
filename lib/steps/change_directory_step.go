package steps

import (
	"log"
	"os"

	"github.com/Originate/git-town/lib/script"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	Directory string
}

// CreateAbortStep returns the abort step for this step.
func (step ChangeDirectoryStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step ChangeDirectoryStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step ChangeDirectoryStep) CreateUndoStep() Step {
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
