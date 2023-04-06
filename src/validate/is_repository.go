package validate

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v7/src/git"
)

// IsRepository verifies that the given folder contains a Git repository.
// It also navigates to the root directory of that repository.
func IsRepository(run *git.ProdRunner) error {
	if !run.Backend.IsRepository() {
		return errors.New("this is not a Git repository")
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current working directory: %w", err)
	}
	gitRootDirectory, err := run.Backend.RootDirectory()
	if err != nil {
		return err
	}
	if currentDirectory != gitRootDirectory {
		return run.Frontend.NavigateToDir(gitRootDirectory)
	}
	return nil
}
