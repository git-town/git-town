package configfile

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

func Encode(config *configdomain.PartialConfig) (string, error) {
	data := toData(config)
	buffer := strings.Builder{}
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(data)
	return buffer.String(), err
}

func Save(config *configdomain.PartialConfig) error {
	text, err := Encode(config)
	if err != nil {
		return err
	}
	file, err := os.Create(FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(text)
	return err
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
