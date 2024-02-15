package validate

import (
	"errors"
	"fmt"
	"slices"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, config *configdomain.FullConfig, localBranches gitdomain.LocalBranchNames, dialogInputs *components.TestInputs) error {
	mainBranch := config.MainBranch
	if mainBranch.IsEmpty() {
		if backend.ConfigFile != nil {
			return errors.New(messages.ConfigMainbranchInConfigFile)
		}
		fmt.Print(messages.ConfigNeeded)
		var err error
		newMainBranch, aborted, err := dialog.MainBranch(localBranches, backend.DefaultBranch(), dialogInputs.Next())
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
		newPerennialBranches, aborted, err := dialog.PerennialBranches(localBranches, config.PerennialBranches, config.MainBranch, dialogInputs.Next())
		if err != nil || aborted {
			return err
		}
		if slices.Compare(newPerennialBranches, config.PerennialBranches) != 0 {
			err := backend.SetPerennialBranches(newPerennialBranches)
			if err != nil {
				return err
			}
		}
	}
	return backend.RemoveOutdatedConfiguration(localBranches)
}
