package validate

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, branches domain.Branches) (domain.BranchTypes, error) {
	mainBranch := backend.GitTown.MainBranch()
	if mainBranch.IsEmpty() {
		fmt.Print("Git Town needs to be configured\n\n")
		newMainBranch, err := dialog.EnterMainBranch(branches.All.LocalBranches().Names(), mainBranch, backend)
		if err != nil {
			return branches.Types, err
		}
		branches.Types.MainBranch = newMainBranch
		return dialog.EnterPerennialBranches(backend, branches)
	}
	return branches.Types, backend.RemoveOutdatedConfiguration(branches.All)
}
