package configfile

import "github.com/git-town/git-town/v11/src/config/configdomain"

//

func Save(config *configdomain.PartialConfig) error {
	data := toData(config)
}

func toData(config *configdomain.PartialConfig) Data {
	result := Data{}
	// branches
	branches := Branches{}
	if config.MainBranch != nil {
		branches.Main = (*string)(config.MainBranch)
	}
	if config.PerennialBranches != nil {
		branches.Perennials = config.PerennialBranches.Strings()
	}
	if !branches.IsEmpty() {
		result.Branches = branches
	}
	// codehosting
	codeHosting := CodeHosting{}
	if config.CodeHostingOriginHostname != nil {
		codeHosting.OriginHostname = (*string)(config.CodeHostingOriginHostname)
	}
	if config.CodeHostingPlatformName != nil {
		codeHosting.Platform = (*string)(config.CodeHostingPlatformName)
	}
	if !codeHosting.IsEmpty() {
		result.CodeHosting = &codeHosting
	}
	// sync-strategy
	syncStrategy := SyncStrategy{}
	if config.SyncFeatureStrategy != nil {
		syncStrategy.FeatureBranches = &config.SyncFeatureStrategy.Name
	}
	if config.SyncPerennialStrategy != nil {
		syncStrategy.PerennialBranches = &config.SyncPerennialStrategy.Name
	}
	if !syncStrategy.IsEmpty() {
		result.SyncStrategy = &syncStrategy
	}
	if config.NewBranchPush != nil {
		result.PushNewbranches = (*bool)(config.NewBranchPush)
	}

	return Data{
		Branches: Branches{Main: (*string)(config.MainBranch.String(), Perennials, []string{}, BadExpr)},
		CodeHosting: &CodeHosting{
			Platform:       new(string),
			OriginHostname: new(string),
		},
		SyncStrategy: &SyncStrategy{
			FeatureBranches:   new(string),
			PerennialBranches: new(string),
		},
		PushNewbranches:          new(bool),
		ShipDeleteTrackingBranch: new(bool),
		SyncUpstream:             new(bool),
	}
}
