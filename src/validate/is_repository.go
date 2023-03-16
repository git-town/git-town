package validate

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
)

// IsRepository verifies that the given folder contains a Git repository.
// It also navigates to the root directory of that repository.
func IsRepository(repo *git.PublicRepo) error {
	if !repo.Silent.IsRepository() {
		return errors.New("this is not a Git repository")
	}
	return repo.NavigateToRootIfNecessary()
}
