package configfile

import (
	"os"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

func Save(config *configdomain.PartialConfig) error {
	text, err := Encode(config)
	if err != nil {
		return err
	}
	return os.WriteFile(FileName, []byte(text), 0o600)
}

func toData(config *configdomain.PartialConfig) Data {
	result := Data{} //nolint:exhaustruct
	// branches
	branches := Branches{} //nolint:exhaustruct
	if config.MainBranch != nil {
		branches.Main = (*string)(config.MainBranch)
	}
	if config.PerennialBranches != nil {
		branches.Perennials = config.PerennialBranches.Strings()
	}
	if !branches.IsEmpty() {
		result.Branches = &branches
	}
	// codehosting
	codeHosting := Hosting{} //nolint:exhaustruct
	if config.HostingOriginHostname != nil {
		codeHosting.OriginHostname = (*string)(config.HostingOriginHostname)
	}
	if config.HostingPlatform != nil {
		codeHosting.Platform = (*string)(config.HostingPlatform)
	}
	if !codeHosting.IsEmpty() {
		result.Hosting = &codeHosting
	}
	// sync-strategy
	syncStrategy := SyncStrategy{} //nolint:exhaustruct
	if config.SyncFeatureStrategy != nil {
		syncStrategy.FeatureBranches = config.SyncFeatureStrategy.StringRef()
	}
	if config.SyncPerennialStrategy != nil {
		syncStrategy.PerennialBranches = config.SyncPerennialStrategy.StringRef()
	}
	if !syncStrategy.IsEmpty() {
		result.SyncStrategy = &syncStrategy
	}
	// top-level fields
	if config.NewBranchPush != nil {
		result.PushNewbranches = (*bool)(config.NewBranchPush)
	}
	if config.PushHook != nil {
		result.PushHook = (*bool)(config.PushHook)
	}
	if config.ShipDeleteTrackingBranch != nil {
		result.ShipDeleteTrackingBranch = (*bool)(config.ShipDeleteTrackingBranch)
	}
	if config.SyncBeforeShip != nil {
		result.SyncBeforeShip = (*bool)(config.SyncBeforeShip)
	}
	if config.SyncUpstream != nil {
		result.SyncUpstream = (*bool)(config.SyncUpstream)
	}
	return result
}
