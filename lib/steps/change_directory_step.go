package steps

import (
	"log"
	"os"

	"github.com/Originate/git-town/lib/script"
)

type ChangeDirectoryStep struct {
	Directory string
}

func (step ChangeDirectoryStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step ChangeDirectoryStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step ChangeDirectoryStep) CreateUndoStep() Step {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return ChangeDirectoryStep{Directory: dir}
}

func (step ChangeDirectoryStep) Run() error {
	_, err := os.Stat(step.Directory)
	if err == nil {
		script.PrintCommand("cd", step.Directory)
		return os.Chdir(step.Directory)
	}
	return nil
}
