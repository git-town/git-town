package configfile

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

// func Save(config *configdomain.PartialConfig) error {
// 	file, err := os.Create(FileName)
// 	if err != nil {
// 		return err
// 	}
// }

func Encode(config *configdomain.PartialConfig) string {
	data := toData(config)
	buffer := strings.Builder{}
	encoder := toml.NewEncoder(&buffer)
	encoder.Encode(data)
	return buffer.String()
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
