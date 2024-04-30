package validate

import (
	"errors"

	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

func ValidateConfig(unvalidated config.UnvalidatedConfig) (*config.ValidatedConfig, error) {
	validateResult, err := MainAndPerennials(MainAndPerennialsArgs{
		UnvalidatedMain:       unvalidated.Config.MainBranch,
		UnvalidatedPerennials: unvalidated.Config.PerennialBranches,
	})
	if err != nil {
		return nil, err
	}
	validatedGitUserEmail, hasGitUserEmail := unvalidated.Config.GitUserEmail.Get()
	if !hasGitUserEmail {
		return nil, errors.New(messages.GitUserEmailMissing)
	}
	validatedGitUserName, hasGitUserName := unvalidated.Config.GitUserName.Get()
	if !hasGitUserName {
		return nil, errors.New(messages.GitUserNameMissing)
	}
	validatedLineage := validateLineage(unvalidated.Config.Lineage)
	validatedConfig := configdomain.ValidatedConfig{
		Aliases:                  unvalidated.Config.Aliases,
		ContributionBranches:     unvalidated.Config.ContributionBranches,
		GitHubToken:              unvalidated.Config.GitHubToken,
		GitLabToken:              unvalidated.Config.GitLabToken,
		GitUserEmail:             validatedGitUserEmail,
		GitUserName:              validatedGitUserName,
		GiteaToken:               unvalidated.Config.GiteaToken,
		HostingOriginHostname:    unvalidated.Config.HostingOriginHostname,
		HostingPlatform:          unvalidated.Config.HostingPlatform,
		Lineage:                  validatedLineage,
		MainBranch:               validateResult.ValidatedMain,
		ObservedBranches:         unvalidated.Config.ObservedBranches,
		Offline:                  unvalidated.Config.Offline,
		ParkedBranches:           unvalidated.Config.ParkedBranches,
		PerennialBranches:        validateResult.ValidatedPerennials,
		PerennialRegex:           unvalidated.Config.PerennialRegex,
		PushHook:                 unvalidated.Config.PushHook,
		PushNewBranches:          unvalidated.Config.PushNewBranches,
		ShipDeleteTrackingBranch: unvalidated.Config.ShipDeleteTrackingBranch,
		SyncBeforeShip:           unvalidated.Config.SyncBeforeShip,
		SyncFeatureStrategy:      unvalidated.Config.SyncFeatureStrategy,
		SyncPerennialStrategy:    unvalidated.Config.SyncPerennialStrategy,
		SyncUpstream:             unvalidated.Config.SyncUpstream,
	}
	vConfig := config.ValidatedConfig{
		ConfigFile:      unvalidated.ConfigFile,
		DryRun:          unvalidated.DryRun,
		FullConfig:      validatedConfig,
		GitConfig:       unvalidated.GitConfig,
		GlobalGitConfig: unvalidated.GlobalGitConfig,
		LocalGitConfig:  unvalidated.LocalGitConfig,
	}
	return vConfig, nil
}
