package validate

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, config *configdomain.FullConfig, localBranches gitdomain.LocalBranchNames) error {
	mainBranch := config.MainBranch
	if mainBranch.IsEmpty() {
		// TODO: extract text
		fmt.Print("Git Town needs to be configured\n\n")
		var err error
		newMainBranch, aborted, err := dialog.EnterMainBranch(localBranches, mainBranch)
		if err != nil || aborted {
			return err
		}
		if newMainBranch != config.MainBranch {
			err := backend.SetMainBranch(newMainBranch)
			if err != nil {
				return err
			}
			config.MainBranch = newMainBranch
		}
		return dialog.EnterPerennialBranches(backend, config, localBranches)
	}
	return backend.RemoveOutdatedConfiguration(localBranches)
}
