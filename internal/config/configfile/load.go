package configfile

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	// keep-sorted start
	var autoResolve Option[configdomain.AutoResolve]
	var autoSync Option[configdomain.AutoSync]
	var branchPrefix Option[configdomain.BranchPrefix]
	var contributionRegex Option[configdomain.ContributionRegex]
	var detached Option[configdomain.Detached]
	var devRemote Option[gitdomain.Remote]
	var displayTypes Option[configdomain.DisplayTypes]
	var featureRegex Option[configdomain.FeatureRegex]
	var forgeType Option[forgedomain.ForgeType]
	var githubConnectorType Option[forgedomain.GitHubConnectorType]
	var gitlabConnectorType Option[forgedomain.GitLabConnectorType]
	var hostingOriginHostname Option[configdomain.HostingOriginHostname]
	var mainBranch Option[gitdomain.LocalBranchName]
	var newBranchType Option[configdomain.NewBranchType]
	var observedRegex Option[configdomain.ObservedRegex]
	var order Option[configdomain.Order]
	var perennialBranches gitdomain.LocalBranchNames
	var perennialRegex Option[configdomain.PerennialRegex]
	var proposalsShowLineage Option[forgedomain.ProposalsShowLineage]
	var pushBranches Option[configdomain.PushBranches]
	var pushHook Option[configdomain.PushHook]
	var shareNewBranches Option[configdomain.ShareNewBranches]
	var shipDeleteTrackingBranch Option[configdomain.ShipDeleteTrackingBranch]
	var shipStrategy Option[configdomain.ShipStrategy]
	var stash Option[configdomain.Stash]
	var syncFeatureStrategy Option[configdomain.SyncFeatureStrategy]
	var syncPerennialStrategy Option[configdomain.SyncPerennialStrategy]
	var syncPrototypeStrategy Option[configdomain.SyncPrototypeStrategy]
	var syncTags Option[configdomain.SyncTags]
	var syncUpstream Option[configdomain.SyncUpstream]
	var unknownBranchType Option[configdomain.UnknownBranchType]
	// keep-sorted end
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
	// keep-sorted start
	var branchPrefixErr error
	var contributionRegexErr error
	var defaultBranchTypeError error
	var displayTypesErr error
	var featureRegexErr error
	var forgeTypeErr error
	var githubConnectorTypeErr error
	var gitlabConnectorTypeErr error
	var hostingOriginHostnameErr error
	var newBranchTypeError error
	var observedRegexErr error
	var orderErr error
	var perennialRegexErr error
	var proposalsShowLineageErr error
	var shareNewBranchesErr error
	var shipDeleteTrackingBranchErr error
	var shipStrategyErr error
	var syncFeatureStrategyErr error
	var syncPerennialStrategyErr error
	var syncPrototypeStrategyErr error
	var syncTagsErr error
	var syncUpstreamErr error
	var unknownBranchTypeError error
	// keep-sorted end
	// load proper definitions, overriding the values from the legacy definitions that were loaded above
	if data.Branches != nil {
		if data.Branches.Main != nil {
			mainBranch = gitdomain.NewLocalBranchNameOption(*data.Branches.Main)
		}
		perennialBranches = gitdomain.NewLocalBranchNames(data.Branches.Perennials...)
		if data.Branches.PerennialRegex != nil {
			perennialRegex, perennialRegexErr = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex, messages.ConfigFile)
		}
		if data.Branches.DefaultType != nil {
			var branchType Option[configdomain.BranchType]
			branchType, defaultBranchTypeError = configdomain.ParseBranchType(*data.Branches.DefaultType, messages.ConfigFile)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
		if data.Branches.DisplayTypes != nil {
			displayTypes, displayTypesErr = configdomain.ParseDisplayTypes(*data.Branches.DisplayTypes, messages.ConfigFile)
		}
		if data.Branches.FeatureRegex != nil {
			var verifiedRegexOpt Option[configdomain.VerifiedRegex]
			verifiedRegexOpt, featureRegexErr = configdomain.ParseRegex(*data.Branches.FeatureRegex)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				featureRegex = Some(configdomain.FeatureRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ContributionRegex != nil {
			var verifiedRegexOpt Option[configdomain.VerifiedRegex]
			verifiedRegexOpt, contributionRegexErr = configdomain.ParseRegex(*data.Branches.ContributionRegex)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				contributionRegex = Some(configdomain.ContributionRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.ObservedRegex != nil {
			var verifiedRegexOpt Option[configdomain.VerifiedRegex]
			verifiedRegexOpt, observedRegexErr = configdomain.ParseRegex(*data.Branches.ObservedRegex)
			if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
				observedRegex = Some(configdomain.ObservedRegex{VerifiedRegex: verifiedRegex})
			}
		}
		if data.Branches.Order != nil {
			order, orderErr = configdomain.ParseOrder(*data.Branches.Order, messages.ConfigFile)
		}
		if data.Branches.UnknownType != nil {
			var branchType Option[configdomain.BranchType]
			branchType, unknownBranchTypeError = configdomain.ParseBranchType(*data.Branches.UnknownType, messages.ConfigFile)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
	}
	if data.Create != nil {
		if data.Create.BranchPrefix != nil {
			branchPrefix, branchPrefixErr = configdomain.ParseBranchPrefix(*data.Create.BranchPrefix, messages.ConfigFile)
		}
		if data.Create.NewBranchType != nil {
			var branchType Option[configdomain.BranchType]
			branchType, newBranchTypeError = configdomain.ParseBranchType(*data.Create.NewBranchType, messages.ConfigFile)
			newBranchType = configdomain.NewBranchTypeOpt(branchType)
		}
		if data.Create.PushNewbranches != nil {
			shareNewBranches = Some(configdomain.ParseShareNewBranchesDeprecatedBool(*data.Create.PushNewbranches))
			finalMessages.Add(messages.PushNewBranchesDeprecation)
		}
		if data.Create.ShareNewBranches != nil {
			shareNewBranches, shareNewBranchesErr = configdomain.ParseShareNewBranches(*data.Create.ShareNewBranches, configdomain.KeyShareNewBranches.String())
		}
		if data.Create.Stash != nil {
			stash = Some(configdomain.Stash(*data.Create.Stash))
		}
	}
	if data.Hosting != nil {
		if data.Hosting.Platform != nil {
			forgeType, forgeTypeErr = forgedomain.ParseForgeType(*data.Hosting.Platform, messages.ConfigFile)
		}
		if data.Hosting.DevRemote != nil {
			devRemote = gitdomain.NewRemote(*data.Hosting.DevRemote)
		}
		if data.Hosting.ForgeType != nil {
			forgeType, forgeTypeErr = forgedomain.ParseForgeType(*data.Hosting.ForgeType, messages.ConfigFile)
		}
		if data.Hosting.GitHubConnectorType != nil {
			githubConnectorType, githubConnectorTypeErr = forgedomain.ParseGitHubConnectorType(*data.Hosting.GitHubConnectorType, messages.ConfigFile)
		}
		if data.Hosting.GitLabConnectorType != nil {
			gitlabConnectorType, gitlabConnectorTypeErr = forgedomain.ParseGitLabConnectorType(*data.Hosting.GitLabConnectorType, messages.ConfigFile)
		}
		if data.Hosting.OriginHostname != nil {
			hostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	if data.Propose != nil {
		if data.Propose.Lineage != nil {
			proposalsShowLineage, proposalsShowLineageErr = forgedomain.ParseProposalsShowLineage(*data.Propose.Lineage, messages.ConfigFile)
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
			syncFeatureStrategy, syncFeatureStrategyErr = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches, messages.ConfigFile)
		}
		if data.SyncStrategy.PerennialBranches != nil {
			syncPerennialStrategy, syncPerennialStrategyErr = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches, messages.ConfigFile)
		}
		if data.SyncStrategy.PrototypeBranches != nil {
			syncPrototypeStrategy, syncPrototypeStrategyErr = configdomain.ParseSyncPrototypeStrategy(*data.SyncStrategy.PrototypeBranches, messages.ConfigFile)
		}
	}
	if data.Sync != nil {
		if data.Sync.AutoResolve != nil {
			autoResolve = Some(configdomain.AutoResolve(*data.Sync.AutoResolve))
		}
		if data.Sync.AutoSync != nil {
			autoSync = Some(configdomain.AutoSync(*data.Sync.AutoSync))
		}
		if data.Sync.Detached != nil {
			detached = Some(configdomain.Detached(*data.Sync.Detached))
		}
		if data.Sync.FeatureStrategy != nil {
			syncFeatureStrategy, syncFeatureStrategyErr = configdomain.ParseSyncFeatureStrategy(*data.Sync.FeatureStrategy, messages.ConfigFile)
		}
		if data.Sync.PerennialStrategy != nil {
			syncPerennialStrategy, syncPerennialStrategyErr = configdomain.ParseSyncPerennialStrategy(*data.Sync.PerennialStrategy, messages.ConfigFile)
		}
		if data.Sync.PrototypeStrategy != nil {
			syncPrototypeStrategy, syncPrototypeStrategyErr = configdomain.ParseSyncPrototypeStrategy(*data.Sync.PrototypeStrategy, messages.ConfigFile)
		}
		if data.Sync.PushBranches != nil {
			pushBranches = Some(configdomain.PushBranches(*data.Sync.PushBranches))
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
	// keep-sorted start
	err := cmp.Or(
		branchPrefixErr,
		defaultBranchTypeError,
		displayTypesErr,
		featureRegexErr,
		contributionRegexErr,
		observedRegexErr,
		perennialRegexErr,
		orderErr,
		unknownBranchTypeError,
		newBranchTypeError,
		shareNewBranchesErr,
		forgeTypeErr,
		githubConnectorTypeErr,
		gitlabConnectorTypeErr,
		hostingOriginHostnameErr,
		proposalsShowLineageErr,
		shipDeleteTrackingBranchErr,
		shipStrategyErr,
		syncFeatureStrategyErr,
		syncPerennialStrategyErr,
		syncPrototypeStrategyErr,
		syncTagsErr,
		syncUpstreamErr,
	// keep-sorted end
	)
	return configdomain.PartialConfig{
		Aliases:                  map[configdomain.AliasableCommand]string{},
		AutoSync:                 autoSync,
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchPrefix:             branchPrefix,
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		ForgejoToken:             None[forgedomain.ForgejoToken](),
		ContributionRegex:        contributionRegex,
		Detached:                 detached,
		DisplayTypes:             displayTypes,
		DryRun:                   None[configdomain.DryRun](),
		UnknownBranchType:        unknownBranchType,
		DevRemote:                devRemote,
		FeatureRegex:             featureRegex,
		ForgeType:                forgeType,
		GitHubConnectorType:      githubConnectorType,
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      gitlabConnectorType,
		GitLabToken:              None[forgedomain.GitLabToken](),
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    hostingOriginHostname,
		Lineage:                  configdomain.NewLineage(),
		MainBranch:               mainBranch,
		NewBranchType:            newBranchType,
		AutoResolve:              autoResolve,
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](),
		Order:                    order,
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		ProposalsShowLineage:     proposalsShowLineage,
		PushBranches:             pushBranches,
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		Stash:                    stash,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		Verbose:                  None[configdomain.Verbose](),
	}, err
}
