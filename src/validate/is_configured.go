package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, allBranches git.BranchesSyncStatus, branchDurations config.BranchDurations) (config.BranchDurations, error) {
	mainBranch := backend.Config.MainBranch()
	if mainBranch == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		newMainBranch, err := EnterMainBranch(allBranches.LocalBranches().BranchNames(), mainBranch, backend)
		if err != nil {
			return branchDurations, err
		}
		branchDurations.MainBranch = newMainBranch
		return EnterPerennialBranches(backend, allBranches, branchDurations)
	}
	return branchDurations, backend.RemoveOutdatedConfiguration(allBranches)
}
