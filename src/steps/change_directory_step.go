package steps

import (
	"os"

	"github.com/git-town/git-town/src/script"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	NoOpStep
	Directory string

	previousDirectory string
}

// CreateUndoStep returns the undo step for this step.
func (step *ChangeDirectoryStep) CreateUndoStep() Step {
	return &ChangeDirectoryStep{Directory: step.previousDirectory}
}

// Run executes this step.
func (step *ChangeDirectoryStep) Run() error {
	var err error
	step.previousDirectory, err = os.Getwd()
	if err != nil {
		return err
	}
	_, err = os.Stat(step.Directory)
	if err == nil {
		script.PrintCommand("cd", step.Directory)
		return os.Chdir(step.Directory)
	}
	return nil
}
