package validate

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, config *configdomain.FullConfig, branches configdomain.Branches) error {
	mainBranch := backend.Config.MainBranch
	if mainBranch.IsEmpty() {
		// TODO: extract text
		fmt.Print("Git Town needs to be configured\n\n")
		var err error
		config.MainBranch, err = dialog.EnterMainBranch(branches.All.LocalBranches().Names(), mainBranch, backend)
		if err != nil {
			return err
		}
		return dialog.EnterPerennialBranches(backend, config, branches)
	}
	return backend.RemoveOutdatedConfiguration(branches.All)
}
