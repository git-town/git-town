package configfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
)

// Decode converts the given config file TOML source into Go data.
func Decode(text string) (*Data, error) {
	var result Data
	_, err := toml.Decode(text, &result)
	return &result, err
}

func Load(rootDir gitdomain.RepoRootDir) (Option[configdomain.PartialConfig], error) {
	configPath := filepath.Join(rootDir.String(), FileName)
	file, err := os.Open(configPath)
	if err != nil {
		return None[configdomain.PartialConfig](), nil //nolint:nilerr
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return None[configdomain.PartialConfig](), fmt.Errorf(messages.ConfigFileCannotRead, ".git-branches.yml", err)
	}
	configFileData, err := Decode(string(bytes))
	if err != nil {
		return None[configdomain.PartialConfig](), fmt.Errorf(messages.ConfigFileInvalidContent, ".git-branches.yml", err)
	}
	result, err := Validate(*configFileData)
	return Some(result), err
}

// Validate converts the given low-level configfile data into high-level config data.
func Validate(data Data) (configdomain.PartialConfig, error) {
	result := configdomain.PartialConfig{} //exhaustruct:ignore
	var err error
	if data.Branches != nil {
		if data.Branches.Main != nil {
			result.MainBranch = gitdomain.NewLocalBranchNameOption(*data.Branches.Main)
		}
		result.PerennialBranches = gitdomain.NewLocalBranchNames(data.Branches.Perennials...)
		if data.Branches.PerennialRegex != nil {
			result.PerennialRegex = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex)
		}
	}
	if data.Hosting != nil {
		if data.Hosting.Platform != nil {
			result.HostingPlatform, err = configdomain.ParseHostingPlatform(*data.Hosting.Platform)
		}
		if data.Hosting.OriginHostname != nil {
			result.HostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	if data.SyncStrategy != nil {
		if data.SyncStrategy.FeatureBranches != nil {
			result.SyncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches)
		}
		if data.SyncStrategy.PerennialBranches != nil {
			result.SyncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches)
		}
	}
	if data.PushNewbranches != nil {
		result.PushNewBranches = Some(configdomain.PushNewBranches(*data.PushNewbranches))
	}
	if data.ShipDeleteTrackingBranch != nil {
		result.ShipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.ShipDeleteTrackingBranch))
	}
	if data.SyncTags != nil {
		result.SyncTags = Some(configdomain.SyncTags(*data.SyncTags))
	}
	if data.SyncUpstream != nil {
		result.SyncUpstream = Some(configdomain.SyncUpstream(*data.SyncUpstream))
	}
	return result, err
}
