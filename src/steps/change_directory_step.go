package steps

import (
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/script"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	NoOpStep
	Directory string

	previousDirectory string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *ChangeDirectoryStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&ChangeDirectoryStep{Directory: step.previousDirectory})
}

// Run executes this step.
func (step *ChangeDirectoryStep) Run() error {
	var err error
	step.previousDirectory, err = os.Getwd()
	exit.If(err)
	_, err = os.Stat(step.Directory)
	if err == nil {
		script.PrintCommand("cd", step.Directory)
		return os.Chdir(step.Directory)
	}
	return nil
}
