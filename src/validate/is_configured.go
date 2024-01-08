package validate

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, config *configdomain.FullConfig, allBranches gitdomain.BranchInfos) error {
	mainBranch := backend.Config.MainBranch
	if mainBranch.IsEmpty() {
		// TODO: extract text
		fmt.Print("Git Town needs to be configured\n\n")
		var err error
		newMainBranch, abort, err := dialog.EnterMainBranch(allBranches.LocalBranches().Names(), mainBranch)
		if err != nil {
			return err
		}
		if abort {
			return nil
		}
		if newMainBranch != config.MainBranch {
			backend.SetMainBranch(config.MainBranch)
		}
		return dialog.EnterPerennialBranches(backend, config, allBranches)
	}
	return backend.RemoveOutdatedConfiguration(allBranches)
}
