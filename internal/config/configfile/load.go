package configfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Decode converts the given config file TOML source into Go data.
func Decode(text string) (*Data, error) {
	var result Data
	_, err := toml.Decode(text, &result)
	return &result, err
}

func Load(rootDir gitdomain.RepoRootDir, fileName string, finalMessages stringslice.Collector) (configdomain.PartialConfig, bool, error) {
	configPath := filepath.Join(rootDir.String(), fileName)
	file, err := os.Open(configPath)
	if err != nil {
		return configdomain.EmptyPartialConfig(), false, nil
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return configdomain.EmptyPartialConfig(), false, fmt.Errorf(messages.ConfigFileCannotRead, fileName, err)
	}
	configFileData, err := Decode(string(bytes))
	if err != nil {
		return configdomain.EmptyPartialConfig(), false, fmt.Errorf(messages.ConfigFileInvalidContent, fileName, err)
	}
	result, err := Validate(*configFileData, finalMessages)
	return result, true, err
}

// Validate converts the given low-level configfile data into high-level config data.
func Validate(data Data, finalMessages stringslice.Collector) (configdomain.PartialConfig, error) {
	var err error
	var contributionRegex Option[configdomain.ContributionRegex]
	var unknownBranchType Option[configdomain.UnknownBranchType]
	var devRemote Option[gitdomain.Remote]
	var featureRegex Option[configdomain.FeatureRegex]
	var forgeType Option[forgedomain.ForgeType]
	var githubConnectorType Option[forgedomain.GitHubConnectorType]
	var gitLabConnectorType Option[forgedomain.GitLabConnectorType]
	var hostingOriginHostname Option[configdomain.HostingOriginHostname]
	var mainBranch Option[gitdomain.LocalBranchName]
	var newBranchType Option[configdomain.NewBranchType]
	var autoResolve Option[configdomain.AutoResolve]
	var observedRegex Option[configdomain.ObservedRegex]
	var perennialBranches gitdomain.LocalBranchNames
	var perennialRegex Option[configdomain.PerennialRegex]
	var proposalsShowLineage Option[forgedomain.ProposalsShowLineage]
	var pushHook Option[configdomain.PushHook]
	var shareNewBranches Option[configdomain.ShareNewBranches]
	var shipDeleteTrackingBranch Option[configdomain.ShipDeleteTrackingBranch]
	var shipStrategy Option[configdomain.ShipStrategy]
	var syncFeatureStrategy Option[configdomain.SyncFeatureStrategy]
	var syncPerennialStrategy Option[configdomain.SyncPerennialStrategy]
	var syncPrototypeStrategy Option[configdomain.SyncPrototypeStrategy]
	var syncTags Option[configdomain.SyncTags]
	var syncUpstream Option[configdomain.SyncUpstream]
	// load legacy definitions first, so that the proper definitions loaded later override them
	if data.CreatePrototypeBranches != nil {
		newBranchType = Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch))
		finalMessages.Add(messages.CreatePrototypeBranchesDeprecation)
	}
	if data.PushNewbranches != nil {
		shareNewBranches = Some(configdomain.ParseShareNewBranchesDeprecatedBool(*data.PushNewbranches))
		finalMessages.Add(messages.PushNewBranchesDeprecation)
	}
	if data.PushHook != nil {
		pushHook = Some(configdomain.PushHook(*data.PushHook))
	}
	if data.ShipDeleteTrackingBranch != nil {
		shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.ShipDeleteTrackingBranch))
	}
	if data.ShipStrategy != nil {
		shipStrategy = Some(configdomain.ShipStrategy(*data.ShipStrategy))
	}
	if data.SyncTags != nil {
		syncTags = Some(configdomain.SyncTags(*data.SyncTags))
	}
	if data.SyncUpstream != nil {
		syncUpstream = Some(configdomain.SyncUpstream(*data.SyncUpstream))
	}
	ec := gohacks.ErrorCollector{}
	// load proper definitions, overriding the values from the legacy definitions that were loaded above
	if data.Branches != nil {
		if data.Branches.Main != nil {
			mainBranch = gitdomain.NewLocalBranchNameOption(*data.Branches.Main)
		}
		perennialBranches = gitdomain.NewLocalBranchNames(data.Branches.Perennials...)
		if data.Branches.PerennialRegex != nil {
			perennialRegex, err = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex)
			ec.Check(err)
		}
		if data.Branches.DefaultType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Branches.DefaultType)
			ec.Check(err)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
		if data.Branches.FeatureRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.FeatureRegex)
			ec.Check(err)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				featureRegex = Some(configdomain.FeatureRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ContributionRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ContributionRegex)
			ec.Check(err)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				contributionRegex = Some(configdomain.ContributionRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ObservedRegex != nil {
			verifiedRegexOpt, err := configdomain.ParseRegex(*data.Branches.ObservedRegex)
			ec.Check(err)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				observedRegex = Some(configdomain.ObservedRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.UnknownType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Branches.UnknownType)
			ec.Check(err)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
	}
	if data.Create != nil {
		if data.Create.NewBranchType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Create.NewBranchType)
			ec.Check(err)
			newBranchType = configdomain.NewBranchTypeOpt(branchType)
		}
		if data.Create.PushNewbranches != nil {
			shareNewBranches = Some(configdomain.ParseShareNewBranchesDeprecatedBool(*data.Create.PushNewbranches))
			finalMessages.Add(messages.PushNewBranchesDeprecation)
		}
		if data.Create.ShareNewBranches != nil {
			shareNewBranches, err = configdomain.ParseShareNewBranches(*data.Create.ShareNewBranches, configdomain.KeyShareNewBranches)
			ec.Check(err)
		}
	}
	if data.Hosting != nil {
		if data.Hosting.Platform != nil {
			forgeType, err = forgedomain.ParseForgeType(*data.Hosting.Platform)
			ec.Check(err)
		}
		if data.Hosting.DevRemote != nil {
			devRemote = gitdomain.NewRemote(*data.Hosting.DevRemote)
		}
		if data.Hosting.ForgeType != nil {
			forgeType, err = forgedomain.ParseForgeType(*data.Hosting.ForgeType)
			ec.Check(err)
		}
		if data.Hosting.GitHubConnectorType != nil {
			githubConnectorType, err = forgedomain.ParseGitHubConnectorType(*data.Hosting.GitHubConnectorType)
			ec.Check(err)
		}
		if data.Hosting.GitLabConnectorType != nil {
			gitLabConnectorType, err = forgedomain.ParseGitLabConnectorType(*data.Hosting.GitLabConnectorType)
			ec.Check(err)
		}
		if data.Hosting.OriginHostname != nil {
			hostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	if data.Propose != nil {
		if data.Propose.Lineage != nil {
			proposalsShowLineage, err = forgedomain.ParseProposalsShowLineage(*data.Propose.Lineage)
			ec.Check(err)
		}
	}
	if data.Ship != nil {
		if data.Ship.DeleteTrackingBranch != nil {
			shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.Ship.DeleteTrackingBranch))
		}
		if data.Ship.Strategy != nil {
			shipStrategy = Some(configdomain.ShipStrategy(*data.Ship.Strategy))
		}
	}
	if data.SyncStrategy != nil {
		if data.SyncStrategy.FeatureBranches != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches)
			ec.Check(err)
		}
		if data.SyncStrategy.PerennialBranches != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches)
			ec.Check(err)
		}
		if data.SyncStrategy.PrototypeBranches != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.SyncStrategy.PrototypeBranches)
			ec.Check(err)
		}
	}
	if data.Sync != nil {
		if data.Sync.AutoResolve != nil {
			autoResolve = Some(configdomain.AutoResolve(*data.Sync.AutoResolve))
		}
		if data.Sync.FeatureStrategy != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.Sync.FeatureStrategy)
			ec.Check(err)
		}
		if data.Sync.PerennialStrategy != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.Sync.PerennialStrategy)
			ec.Check(err)
		}
		if data.Sync.PrototypeStrategy != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.Sync.PrototypeStrategy)
			ec.Check(err)
		}
		if data.Sync.PushHook != nil {
			pushHook = Some(configdomain.PushHook(*data.Sync.PushHook))
		}
		if data.Sync.Tags != nil {
			syncTags = Some(configdomain.SyncTags(*data.Sync.Tags))
		}
		if data.Sync.Upstream != nil {
			syncUpstream = Some(configdomain.SyncUpstream(*data.Sync.Upstream))
		}
	}
	return configdomain.PartialConfig{
		Aliases:                  map[configdomain.AliasableCommand]string{},
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		CodebergToken:            None[forgedomain.CodebergToken](),
		ContributionRegex:        contributionRegex,
		DryRun:                   None[configdomain.DryRun](),
		UnknownBranchType:        unknownBranchType,
		DevRemote:                devRemote,
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      gitLabConnectorType,
		GitLabToken:              None[forgedomain.GitLabToken](),
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    hostingOriginHostname,
		Lineage:                  configdomain.Lineage{},
		MainBranch:               mainBranch,
		NewBranchType:            newBranchType,
		AutoResolve:              autoResolve,
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](),
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		ProposalsShowLineage:     proposalsShowLineage,
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		Verbose:                  None[configdomain.Verbose](),
	}, ec.Err
}
