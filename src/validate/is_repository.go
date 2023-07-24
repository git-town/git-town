package validate

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/git"
)

// IsRepository verifies that the given folder contains a Git repository.
// It also navigates to the root directory of that repository.
func IsRepository(run *git.ProdRunner) error {
	isRepo, repoDir := run.Backend.IsRepository()
	if !isRepo {
		return errors.New("this is not a Git repository")
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current working directory: %w", err)
	}
	if currentDirectory != repoDir {
		return run.Frontend.NavigateToDir(repoDir)
	}
	return nil
}
