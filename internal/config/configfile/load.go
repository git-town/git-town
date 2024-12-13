package configfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Decode converts the given config file TOML source into Go data.
func Decode(text string) (*Data, error) {
	var result Data
	_, err := toml.Decode(text, &result)
	return &result, err
}

func Load(rootDir gitdomain.RepoRootDir, finalMessages stringslice.Collector) (Option[configdomain.PartialConfig], error) {
	configPath := filepath.Join(rootDir.String(), FileName)
	file, err := os.Open(configPath)
	if err != nil {
		return None[configdomain.PartialConfig](), nil
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
	result, err := Validate(*configFileData, finalMessages)
	return Some(result), err
}

// Validate converts the given low-level configfile data into high-level config data.
func Validate(data Data, finalMessages stringslice.Collector) (configdomain.PartialConfig, error) {
	var err error
	var mainBranch Option[gitdomain.LocalBranchName]
	var perennialBranches gitdomain.LocalBranchNames
	var perennialRegex Option[configdomain.PerennialRegex]
	var defaultBranchType Option[configdomain.BranchType]
	var featureRegex Option[configdomain.FeatureRegex]
	var contributionRegex Option[configdomain.ContributionRegex]
	var observedRegex Option[configdomain.ObservedRegex]
	if data.Branches != nil {
		if data.Branches.Main != nil {
			mainBranch = gitdomain.NewLocalBranchNameOption(*data.Branches.Main)
		}
		perennialBranches = gitdomain.NewLocalBranchNames(data.Branches.Perennials...)
		if data.Branches.PerennialRegex != nil {
			perennialRegex, err = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Branches.DefaultType != nil {
			defaultBranchType, err = configdomain.ParseBranchType(*data.Branches.DefaultType)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Branches.FeatureRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.FeatureRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				featureRegex = Some(configdomain.FeatureRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ContributionRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ContributionRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				contributionRegex = Some(configdomain.ContributionRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ObservedRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ObservedRegex)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				observedRegex = Some(configdomain.ObservedRegex{VerifiedRegex: verifiedRegex})
			}
		}
	}
	var newBranchType Option[configdomain.BranchType]
	if data.CreatePrototypeBranches != nil {
		newBranchType = Some(configdomain.BranchTypePrototypeBranch)
		finalMessages.Add(messages.CreatePrototypeBranchesDeprecation)
	}
	if data.Create != nil {
		if data.Create.NewBranchType != nil {
			parsed, err := configdomain.ParseBranchType(*data.Create.NewBranchType)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
			newBranchType = parsed
		}
	}
	var hostingPlatform Option[configdomain.HostingPlatform]
	var hostingOriginHostname Option[configdomain.HostingOriginHostname]
	if data.Hosting != nil {
		if data.Hosting.Platform != nil {
			hostingPlatform, err = configdomain.ParseHostingPlatform(*data.Hosting.Platform)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.Hosting.OriginHostname != nil {
			hostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	var syncFeatureStrategy Option[configdomain.SyncFeatureStrategy]
	var syncPerennialStrategy Option[configdomain.SyncPerennialStrategy]
	var syncPrototypeStrategy Option[configdomain.SyncPrototypeStrategy]
	if data.SyncStrategy != nil {
		if data.SyncStrategy.FeatureBranches != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.SyncStrategy.PerennialBranches != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
		if data.SyncStrategy.PrototypeBranches != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.SyncStrategy.PrototypeBranches)
			if err != nil {
				return configdomain.EmptyPartialConfig(), err
			}
		}
	}
	var pushNewBranches Option[configdomain.PushNewBranches]
	if data.PushNewbranches != nil {
		pushNewBranches = Some(configdomain.PushNewBranches(*data.PushNewbranches))
	}
	var pushHook Option[configdomain.PushHook]
	if data.PushHook != nil {
		pushHook = Some(configdomain.PushHook(*data.PushHook))
	}
	var shipDeleteTrackingBranch Option[configdomain.ShipDeleteTrackingBranch]
	if data.ShipDeleteTrackingBranch != nil {
		shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.ShipDeleteTrackingBranch))
	}
	var shipStrategy Option[configdomain.ShipStrategy]
	if data.ShipStrategy != nil {
		shipStrategy = Some(configdomain.ShipStrategy(*data.ShipStrategy))
	}
	if data.Ship != nil {
		if data.Ship.DeleteTrackingBranch != nil {
			shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.Ship.DeleteTrackingBranch))
		}
		if data.Ship.Strategy != nil {
			shipStrategy = Some(configdomain.ShipStrategy(*data.Ship.Strategy))
		}
	}
	var syncTags Option[configdomain.SyncTags]
	if data.SyncTags != nil {
		syncTags = Some(configdomain.SyncTags(*data.SyncTags))
	}
	var syncUpstream Option[configdomain.SyncUpstream]
	if data.SyncUpstream != nil {
		syncUpstream = Some(configdomain.SyncUpstream(*data.SyncUpstream))
	}
	if data.Sync != nil {
		if data.Sync.PushHook != nil {
			pushHook = Some(configdomain.PushHook(*data.Sync.PushHook))
		}
	}
	return configdomain.PartialConfig{
		Aliases:                  map[configdomain.AliasableCommand]string{},
		BitbucketAppPassword:     None[configdomain.BitbucketAppPassword](),
		BitbucketUsername:        None[configdomain.BitbucketUsername](),
		ContributionBranches:     gitdomain.LocalBranchNames{},
		ContributionRegex:        contributionRegex,
		DefaultBranchType:        defaultBranchType,
		FeatureRegex:             featureRegex,
		GitHubToken:              None[configdomain.GitHubToken](),
		GitLabToken:              None[configdomain.GitLabToken](),
		GitUserEmail:             None[configdomain.GitUserEmail](),
		GitUserName:              None[configdomain.GitUserName](),
		GiteaToken:               None[configdomain.GiteaToken](),
		HostingOriginHostname:    hostingOriginHostname,
		HostingPlatform:          hostingPlatform,
		Lineage:                  configdomain.Lineage{},
		MainBranch:               mainBranch,
		NewBranchType:            newBranchType,
		ObservedBranches:         gitdomain.LocalBranchNames{},
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](),
		ParkedBranches:           gitdomain.LocalBranchNames{},
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		PrototypeBranches:        gitdomain.LocalBranchNames{},
		PushHook:                 pushHook,
		PushNewBranches:          pushNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
	}, nil
}
