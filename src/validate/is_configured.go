package validate

import (
	"errors"
	"fmt"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(backend *git.BackendCommands, config *config.Config, localBranches gitdomain.LocalBranchNames, dialogInputs *components.TestInputs) error {
	mainBranch := config.Config.MainBranch
	if mainBranch.IsEmpty() {
		if config.ConfigFile.IsSome() {
			return errors.New(messages.ConfigMainbranchInConfigFile)
		}
		fmt.Print(messages.ConfigNeeded)
		var err error
		newMainBranch, aborted, err := dialog.MainBranch(localBranches, backend.DefaultBranch(), dialogInputs.Next())
		if err != nil || aborted {
			return err
		}
		if newMainBranch != config.Config.MainBranch {
			err := config.SetMainBranch(newMainBranch)
			if err != nil {
				return err
			}
			config.Config.MainBranch = newMainBranch
		}
		newPerennialBranches, aborted, err := dialog.PerennialBranches(localBranches, config.Config.PerennialBranches, config.Config.MainBranch, dialogInputs.Next())
		if err != nil || aborted {
			return err
		}
		if slices.Compare(newPerennialBranches, config.Config.PerennialBranches) != 0 {
			err := config.SetPerennialBranches(newPerennialBranches)
			if err != nil {
				return err
			}
		}
	}
	return config.RemoveOutdatedConfiguration(localBranches)
}
