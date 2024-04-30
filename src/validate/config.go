package validate

import "github.com/git-town/git-town/v14/src/config/configdomain"

func ValidateConfig(unvalidated configdomain.UnvalidatedConfig) (*configdomain.ValidatedConfig, error) {
	validateResult, err := MainAndPerennials(MainAndPerennialsArgs{
		UnvalidatedMain:       unvalidated.MainBranch,
		UnvalidatedPerennials: unvalidated.PerennialBranches,
	})
	if err != nil {
		return nil, err
	}
	validatedGitUserEmail := validateGitUserEmail(unvalidated.GitUserEmail)
	validatedGitUserName := validateGitUserName(unvalidated.GitUserName)
	validatedLineage := validateLineage(unvalidated.Lineage)
	return &configdomain.ValidatedConfig{
		Aliases:                  unvalidated.Aliases,
		ContributionBranches:     unvalidated.ContributionBranches,
		GitHubToken:              unvalidated.GitHubToken,
		GitLabToken:              unvalidated.GitLabToken,
		GitUserEmail:             validatedGitUserEmail,
		GitUserName:              validatedGitUserName,
		GiteaToken:               unvalidated.GiteaToken,
		HostingOriginHostname:    unvalidated.HostingOriginHostname,
		HostingPlatform:          unvalidated.HostingPlatform,
		Lineage:                  validatedLineage,
		MainBranch:               validateResult.ValidatedMain,
		ObservedBranches:         unvalidated.ObservedBranches,
		Offline:                  unvalidated.Offline,
		ParkedBranches:           unvalidated.ParkedBranches,
		PerennialBranches:        validateResult.ValidatedPerennials,
		PerennialRegex:           unvalidated.PerennialRegex,
		PushHook:                 unvalidated.PushHook,
		PushNewBranches:          unvalidated.PushNewBranches,
		ShipDeleteTrackingBranch: unvalidated.ShipDeleteTrackingBranch,
		SyncBeforeShip:           unvalidated.SyncBeforeShip,
		SyncFeatureStrategy:      unvalidated.SyncFeatureStrategy,
		SyncPerennialStrategy:    unvalidated.SyncPerennialStrategy,
		SyncUpstream:             unvalidated.SyncUpstream,
	}, nil
}
