package steps

import (
	"os"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// ChangeDirectoryStep changes the current working directory.
type ChangeDirectoryStep struct {
	NoOpStep
	Directory string

	previousDirectory string
}

// CreateUndoStep returns the undo step for this step.
func (step *ChangeDirectoryStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ChangeDirectoryStep{Directory: step.previousDirectory}, nil
}

// Run executes this step.
func (step *ChangeDirectoryStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	var err error
	step.previousDirectory, err = os.Getwd()
	if err != nil {
		return err
	}
	_, err = os.Stat(step.Directory)
	if err == nil {
		err = repo.LoggingShell.PrintCommand("cd", step.Directory)
		if err != nil {
			return err
		}
		return os.Chdir(step.Directory)
	}
	return nil
}
