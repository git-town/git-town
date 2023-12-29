package configfile

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

const ConfigFileName = ".git-branches.toml"

// ConfigFileData is the unvalidated data as read by the TOML parser.
type ConfigFileData struct {
	Branches                 ConfigFileBranches `toml:"branches"`
	CodeHosting              *CodeHosting       `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy      `toml:"sync-strategy"`
	PushNewbranches          *bool              `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool              `toml:"ship-delete-remote-branch"`
	SyncUpstream             *bool              `toml:"sync-upstream"`
}

type ConfigFileBranches struct {
	Main       *string  `toml:"main"`
	Perennials []string `toml:"perennials"`
}

type CodeHosting struct {
	Platform       *string `toml:"platform"`
	OriginHostname *string `toml:"origin-hostname"`
}

type SyncStrategy struct {
	FeatureBranches   *string `toml:"feature-branches"`
	PerennialBranches *string `toml:"perennial-branches"`
}

func (self ConfigFileData) Validate() (configdomain.PartialConfig, error) {
	result := configdomain.PartialConfig{} //nolint:exhaustruct
	var err error
	if self.Branches.Main != nil {
		result.MainBranch = gitdomain.NewLocalBranchNameRef(*self.Branches.Main)
	}
	if self.Branches.Perennials != nil {
		result.PerennialBranches = gitdomain.NewLocalBranchNamesRef(self.Branches.Perennials...)
	}
	if self.CodeHosting != nil {
		if self.CodeHosting.Platform != nil {
			result.CodeHostingPlatformName = configdomain.NewCodeHostingPlatformNameRef(*self.CodeHosting.Platform)
		}
		if self.CodeHosting.OriginHostname != nil {
			result.CodeHostingOriginHostname = configdomain.NewCodeHostingOriginHostnameRef(*self.CodeHosting.OriginHostname)
		}
	}
	if self.SyncStrategy != nil {
		if self.SyncStrategy.FeatureBranches != nil {
			result.SyncFeatureStrategy, err = configdomain.NewSyncFeatureStrategyRef(*self.SyncStrategy.FeatureBranches)
		}
		if self.SyncStrategy.PerennialBranches != nil {
			result.SyncPerennialStrategy, err = configdomain.NewSyncPerennialStrategyRef(*self.SyncStrategy.PerennialBranches)
		}
	}
	if self.PushNewbranches != nil {
		result.NewBranchPush = configdomain.NewNewBranchPushRef(*self.PushNewbranches)
	}
	if self.ShipDeleteTrackingBranch != nil {
		result.ShipDeleteTrackingBranch = configdomain.NewShipDeleteTrackingBranchRef(*self.ShipDeleteTrackingBranch)
	}
	if self.SyncUpstream != nil {
		result.SyncUpstream = configdomain.NewSyncUpstreamRef(*self.SyncUpstream)
	}
	return result, err
}

func LoadConfigFile() (configdomain.PartialConfig, error) {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return configdomain.EmptyPartialConfig(), nil //nolint:nilerr
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return configdomain.EmptyPartialConfig(), fmt.Errorf(messages.ConfigFileCannotRead, ".git-branches.yml", err)
	}
	configFileData, err := ParseTOML(string(bytes))
	if err != nil {
		return configdomain.EmptyPartialConfig(), fmt.Errorf(messages.ConfigFileInvalidData, ".git-branches.yml", err)
	}
	return configFileData.Validate()
}

func ParseTOML(text string) (*ConfigFileData, error) {
	var result ConfigFileData
	_, err := toml.Decode(text, &result)
	return &result, err
}
