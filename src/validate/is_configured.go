package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, allBranches domain.BranchInfos, branchDurations domain.BranchDurations) (domain.BranchDurations, error) {
	mainBranch := backend.Config.MainBranch()
	if mainBranch.IsEmpty() {
		fmt.Print("Git Town needs to be configured\n\n")
		newMainBranch, err := dialog.EnterMainBranch(allBranches.LocalBranches().Names(), mainBranch, backend)
		if err != nil {
			return branchDurations, err
		}
		branchDurations.MainBranch = newMainBranch
		return dialog.EnterPerennialBranches(backend, allBranches, branchDurations)
	}
	return branchDurations, backend.RemoveOutdatedConfiguration(allBranches)
}
