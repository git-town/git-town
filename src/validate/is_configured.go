package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, allBranches git.BranchesSyncStatus) error {
	mainBranch := backend.Config.MainBranch()
	if mainBranch == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		mainBranch, err := EnterMainBranch(mainBranch, backend)
		if err != nil {
			return err
		}
		return EnterPerennialBranches(backend, allBranches, mainBranch)
	}
	return backend.RemoveOutdatedConfiguration(allBranches)
}
