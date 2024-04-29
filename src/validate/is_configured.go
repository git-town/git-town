package validate

import (
	"errors"
	"fmt"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(unvalidatedConfig configdomain.UnvalidatedConfig, backend *git.BackendCommands, config *config.Config, localBranches gitdomain.LocalBranchNames, dialogInputs *components.TestInputs) (*configdomain.ValidatedConfig, error) {
	var mainBranch gitdomain.LocalBranchName
	var perennialBranches gitdomain.LocalBranchNames
	if unvalidatedConfig.MainBranch.IsEmpty() {
		if config.ConfigFile.IsSome() {
			return nil, errors.New(messages.ConfigMainbranchInConfigFile)
		}
		fmt.Print(messages.ConfigNeeded)
		var err error
		var aborted bool
		mainBranch, aborted, err = dialog.MainBranch(localBranches, backend.DefaultBranch(), dialogInputs.Next())
		if err != nil || aborted {
			return nil, err
		}
		if mainBranch != unvalidatedConfig.MainBranch {
			err := config.SetMainBranch(mainBranch)
			if err != nil {
				return nil, err
			}
		}
		perennialBranches, aborted, err = dialog.PerennialBranches(localBranches, unvalidatedConfig.PerennialBranches, config.FullConfig.MainBranch, dialogInputs.Next())
		if err != nil || aborted {
			return nil, err
		}
		if slices.Compare(perennialBranches, config.FullConfig.PerennialBranches) != 0 {
			err := config.SetPerennialBranches(perennialBranches)
			if err != nil {
				return nil, err
			}
		}
	} else {
		mainBranch = unvalidatedConfig.MainBranch
		perennialBranches = unvalidatedConfig.PerennialBranches
	}
	err := config.RemoveOutdatedConfiguration(localBranches)
	return &configdomain.ValidatedConfig{
		Aliases:                  unvalidatedConfig.Aliases,
		ContributionBranches:     unvalidatedConfig.ContributionBranches,
		GitHubToken:              unvalidatedConfig.GitHubToken,
		GitLabToken:              unvalidatedConfig.GitLabToken,
		GitUserEmail:             unvalidatedConfig.GitUserEmail,
		GitUserName:              unvalidatedConfig.GitUserName,
		GiteaToken:               unvalidatedConfig.GiteaToken,
		HostingOriginHostname:    unvalidatedConfig.HostingOriginHostname,
		HostingPlatform:          unvalidatedConfig.HostingPlatform,
		Lineage:                  unvalidatedConfig.Lineage,
		MainBranch:               mainBranch,
		ObservedBranches:         unvalidatedConfig.ObservedBranches,
		Offline:                  unvalidatedConfig.Offline,
		ParkedBranches:           unvalidatedConfig.ParkedBranches,
		PerennialBranches:        perennialBranches,
		PerennialRegex:           unvalidatedConfig.PerennialRegex,
		PushHook:                 unvalidatedConfig.PushHook,
		PushNewBranches:          unvalidatedConfig.PushNewBranches,
		ShipDeleteTrackingBranch: unvalidatedConfig.ShipDeleteTrackingBranch,
		SyncBeforeShip:           unvalidatedConfig.SyncBeforeShip,
		SyncFeatureStrategy:      unvalidatedConfig.SyncFeatureStrategy,
		SyncPerennialStrategy:    unvalidatedConfig.SyncPerennialStrategy,
		SyncUpstream:             unvalidatedConfig.SyncUpstream,
	}, err
}
