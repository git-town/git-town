package validate

func ValidateConfig(unvalidated UnvalidatedConfig) ValidatedConfig {
	validateResult := MainAndPerennials(MainAndPerennialsArgs{
		UnvalidatedMain:       unvalidated.MainBranch,
		UnvalidatedPerennials: unvalidated.PerennialBranches,
	})
	validatedGitUserEmail := validateGitUserEmail(unvalidated.GitUserEmail)
	validatedGitUserName := validateGitUserName(unvalidated.GitUserName)
	validatedLineage := validateLineage(unvalidated.Lineage)
	return ValidatedConfig{
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
		MainBranch:               validatedMainBranch,
		ObservedBranches:         unvalidated.ObservedBranches,
		Offline:                  unvalidated.Offline,
		ParkedBranches:           unvalidated.ParkedBranches,
		PerennialBranches:        validatedPerennialBranches,
		PerennialRegex:           unvalidated.PerennialRegex,
		PushHook:                 unvalidated.PushHook,
		PushNewBranches:          unvalidated.PushNewBranches,
		ShipDeleteTrackingBranch: unvalidated.ShipDeleteTrackingBranch,
		SyncBeforeShip:           unvalidated.SyncBeforeShip,
		SyncFeatureStrategy:      unvalidated.SyncFeatureStrategy,
		SyncPerennialStrategy:    unvalidated.SyncPerennialStrategy,
		SyncUpstream:             unvalidated.SyncUpstream,
	}
}
