package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands) error {
	if backend.Config.MainBranch() == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		mainBranch, err := EnterMainBranch(backend)
		if err != nil {
			return err
		}
		return EnterPerennialBranches(backend, mainBranch)
	}
	return backend.RemoveOutdatedConfiguration()
}
