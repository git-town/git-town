package configdomain

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/messages"
)

// ConfigFileData is the unvalidated data as read by the TOML parser.
type ConfigFileData struct {
	Branches                 Branches      `toml:"branches"`
	CodeHosting              *CodeHosting  `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy `toml:"sync-strategy"`
	PushNewbranches          *bool         `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *bool         `toml:"ship-delete-remote-branch"`
	SyncUpstream             *bool         `toml:"sync-upstream"`
}

type Branches struct {
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

func (self ConfigFileData) Validate() (PartialGitConfig, error) {
	result := PartialGitConfig{} //nolint:exhaustruct
	var err error
	if self.Branches.Main != nil {
		result.MainBranch = domain.NewLocalBranchNameRef(*self.Branches.Main)
	}
	if self.Branches.Perennials != nil {
		result.PerennialBranches = domain.NewLocalBranchNamesRef(self.Branches.Perennials...)
	}
	if self.CodeHosting != nil {
		if self.CodeHosting.Platform != nil {
			result.CodeHostingPlatformName = NewCodeHostingPlatformNameRef(*self.CodeHosting.Platform)
		}
		if self.CodeHosting.OriginHostname != nil {
			result.CodeHostingOriginHostname = NewCodeHostingOriginHostnameRef(*self.CodeHosting.OriginHostname)
		}
	}
	if self.SyncStrategy != nil {
		if self.SyncStrategy.FeatureBranches != nil {
			result.SyncFeatureStrategy, err = NewSyncFeatureStrategyRef(*self.SyncStrategy.FeatureBranches)
		}
		if self.SyncStrategy.PerennialBranches != nil {
			result.SyncPerennialStrategy, err = NewSyncPerennialStrategyRef(*self.SyncStrategy.PerennialBranches)
		}
	}
	if self.PushNewbranches != nil {
		result.NewBranchPush = NewNewBranchPushRef(*self.PushNewbranches)
	}
	if self.ShipDeleteTrackingBranch != nil {
		result.ShipDeleteTrackingBranch = NewShipDeleteTrackingBranchRef(*self.ShipDeleteTrackingBranch)
	}
	if self.SyncUpstream != nil {
		result.SyncUpstream = NewSyncUpstreamRef(*self.SyncUpstream)
	}
	return result, err
}

// ConfigFile is validated data from the configuration file, ready to be used by the application.
type ConfigFile struct {
	Branches                 Branches                  `toml:"branches"`
	CodeHosting              *CodeHosting              `toml:"code-hosting"`
	SyncStrategy             *SyncStrategy             `toml:"sync-strategy"`
	PushNewbranches          *NewBranchPush            `toml:"push-new-branches"`
	ShipDeleteTrackingBranch *ShipDeleteTrackingBranch `toml:"ship-delete-remote-branch"`
	SyncUpstream             *SyncUpstream             `toml:"sync-upstream"`
}

func LoadConfigFile() (PartialGitConfig, error) {
	file, err := os.Open(".git-branches.toml")
	if err != nil {
		return EmptyPartialConfig(), nil //nolint:nilerr
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return EmptyPartialConfig(), fmt.Errorf(messages.ConfigFileCannotRead, ".git-branches.yml", err)
	}
	configFileData, err := ParseTOML(string(bytes))
	if err != nil {
		return EmptyPartialConfig(), fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return configFileData.Validate()
}

func ParseTOML(text string) (*ConfigFileData, error) {
	var result ConfigFileData
	_, err := toml.Decode(text, &result)
	if err != nil {
		return nil, fmt.Errorf(messages.ConfigFileWrongInput, ".git-branches.yml", err)
	}
	return &result, err
}
