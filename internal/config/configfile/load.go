package configfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
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
	// TODO: convert to proper variable initialization using None
	var (
		// keep-sorted start
		autoResolve                 Option[configdomain.AutoResolve]
		autoSync                    Option[configdomain.AutoSync]
		branchPrefix                Option[configdomain.BranchPrefix]
		browser                     Option[configdomain.Browser]
		contributionRegex           Option[configdomain.ContributionRegex]
		detached                    Option[configdomain.Detached]
		devRemote                   Option[gitdomain.Remote]
		displayTypes                Option[configdomain.DisplayTypes]
		featureRegex                Option[configdomain.FeatureRegex]
		forgeType                   Option[forgedomain.ForgeType]
		githubConnectorType         Option[forgedomain.GithubConnectorType]
		gitlabConnectorType         Option[forgedomain.GitlabConnectorType]
		hostingOriginHostname       Option[configdomain.HostingOriginHostname]
		ignoreUncommitted           Option[configdomain.IgnoreUncommitted]
		mainBranch                  Option[gitdomain.LocalBranchName]
		newBranchType               Option[configdomain.NewBranchType]
		observedRegex               Option[configdomain.ObservedRegex]
		order                       Option[configdomain.Order]
		perennialBranches           gitdomain.LocalBranchNames
		perennialRegex              Option[configdomain.PerennialRegex]
		proposalBreadcrumb          Option[configdomain.ProposalBreadcrumb]
		proposalBreadcrumbDirection Option[configdomain.ProposalBreadcrumbDirection]
		pushBranches                Option[configdomain.PushBranches]
		pushHook                    Option[configdomain.PushHook]
		shareNewBranches            Option[configdomain.ShareNewBranches]
		shipDeleteTrackingBranch    Option[configdomain.ShipDeleteTrackingBranch]
		shipStrategy                Option[configdomain.ShipStrategy]
		stash                       Option[configdomain.Stash]
		syncFeatureStrategy         Option[configdomain.SyncFeatureStrategy]
		syncPerennialStrategy       Option[configdomain.SyncPerennialStrategy]
		syncPrototypeStrategy       Option[configdomain.SyncPrototypeStrategy]
		syncTags                    Option[configdomain.SyncTags]
		syncUpstream                Option[configdomain.SyncUpstream]
		unknownBranchType           Option[configdomain.UnknownBranchType]
		// keep-sorted end
	)
	var err error
	// load legacy definitions first, so that the proper definitions loaded later override them
	if data.CreatePrototypeBranches != nil {
		newBranchType = Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch))
		finalMessages.Add(messages.CreatePrototypeBranchesDeprecation)
	}
	if data.PushNewBranches != nil {
		shareNewBranches = Some(configdomain.ParseShareNewBranchesDeprecatedBool(*data.PushNewBranches))
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
			perennialRegex, err = configdomain.ParsePerennialRegex(*data.Branches.PerennialRegex, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Branches.DefaultType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Branches.DefaultType, messages.ConfigFile)
			ec.Check(err)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
		if data.Branches.DisplayTypes != nil {
			displayTypes, err = configdomain.ParseDisplayTypes(*data.Branches.DisplayTypes, messages.ConfigFile)
			ec.Check(err)
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
		if data.Branches.Order != nil {
			order, err = configdomain.ParseOrder(*data.Branches.Order, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Branches.UnknownType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Branches.UnknownType, messages.ConfigFile)
			ec.Check(err)
			unknownBranchType = configdomain.UnknownBranchTypeOpt(branchType)
		}
	}
	if data.Create != nil {
		if data.Create.BranchPrefix != nil {
			branchPrefix, err = configdomain.ParseBranchPrefix(*data.Create.BranchPrefix, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Create.NewBranchType != nil {
			branchType, err := configdomain.ParseBranchType(*data.Create.NewBranchType, messages.ConfigFile)
			ec.Check(err)
			newBranchType = configdomain.NewBranchTypeOpt(branchType)
		}
		if data.Create.PushNewBranches != nil {
			shareNewBranches = Some(configdomain.ParseShareNewBranchesDeprecatedBool(*data.Create.PushNewBranches))
			finalMessages.Add(messages.PushNewBranchesDeprecation)
		}
		if data.Create.ShareNewBranches != nil {
			shareNewBranches, err = configdomain.ParseShareNewBranches(*data.Create.ShareNewBranches, configdomain.KeyShareNewBranches.String())
			ec.Check(err)
		}
		if data.Create.Stash != nil {
			stash = Some(configdomain.Stash(*data.Create.Stash))
		}
	}
	if data.Hosting != nil {
		if data.Hosting.Browser != nil {
			browser, err = configdomain.ParseBrowser(*data.Hosting.Browser, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Hosting.Platform != nil {
			forgeType, err = forgedomain.ParseForgeType(*data.Hosting.Platform, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Hosting.DevRemote != nil {
			devRemote = gitdomain.NewRemote(*data.Hosting.DevRemote)
		}
		if data.Hosting.ForgeType != nil {
			forgeType, err = forgedomain.ParseForgeType(*data.Hosting.ForgeType, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Hosting.GithubConnector != nil {
			githubConnectorType, err = forgedomain.ParseGithubConnectorType(*data.Hosting.GithubConnector, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Hosting.GitlabConnector != nil {
			gitlabConnectorType, err = forgedomain.ParseGitlabConnectorType(*data.Hosting.GitlabConnector, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Hosting.OriginHostname != nil {
			hostingOriginHostname = configdomain.ParseHostingOriginHostname(*data.Hosting.OriginHostname)
		}
	}
	if data.Propose != nil {
		// load the deprecated "lineage" setting first so that "breadcrumb" can override the value later
		if data.Propose.Lineage != nil {
			proposalBreadcrumb, err = configdomain.ParseProposalBreadcrumb(*data.Propose.Lineage, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Propose.Breadcrumb != nil {
			proposalBreadcrumb, err = configdomain.ParseProposalBreadcrumb(*data.Propose.Breadcrumb, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Propose.Direction != nil {
			proposalBreadcrumbDirection, err = configdomain.ParseProposalBreadcrumbDirection(*data.Propose.Direction, messages.ConfigFile)
			ec.Check(err)
		}
	}
	if data.Ship != nil {
		if data.Ship.DeleteTrackingBranch != nil {
			shipDeleteTrackingBranch = Some(configdomain.ShipDeleteTrackingBranch(*data.Ship.DeleteTrackingBranch))
		}
		if data.Ship.IgnoreUncommitted != nil {
			ignoreUncommitted = Some(configdomain.IgnoreUncommitted(*data.Ship.IgnoreUncommitted))
		}
		if data.Ship.Strategy != nil {
			shipStrategy = Some(configdomain.ShipStrategy(*data.Ship.Strategy))
		}
	}
	if data.SyncStrategy != nil {
		if data.SyncStrategy.FeatureBranches != nil {
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.SyncStrategy.FeatureBranches, messages.ConfigFile)
			ec.Check(err)
		}
		if data.SyncStrategy.PerennialBranches != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.SyncStrategy.PerennialBranches, messages.ConfigFile)
			ec.Check(err)
		}
		if data.SyncStrategy.PrototypeBranches != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.SyncStrategy.PrototypeBranches, messages.ConfigFile)
			ec.Check(err)
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
			syncFeatureStrategy, err = configdomain.ParseSyncFeatureStrategy(*data.Sync.FeatureStrategy, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Sync.PerennialStrategy != nil {
			syncPerennialStrategy, err = configdomain.ParseSyncPerennialStrategy(*data.Sync.PerennialStrategy, messages.ConfigFile)
			ec.Check(err)
		}
		if data.Sync.PrototypeStrategy != nil {
			syncPrototypeStrategy, err = configdomain.ParseSyncPrototypeStrategy(*data.Sync.PrototypeStrategy, messages.ConfigFile)
			ec.Check(err)
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
	return configdomain.PartialConfig{
		Aliases:                     map[configdomain.AliasableCommand]string{},
		AutoSync:                    autoSync,
		BitbucketAppPassword:        None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:           None[forgedomain.BitbucketUsername](),
		BranchPrefix:                branchPrefix,
		BranchTypeOverrides:         configdomain.BranchTypeOverrides{},
		Browser:                     browser,
		ForgejoToken:                None[forgedomain.ForgejoToken](),
		ContributionRegex:           contributionRegex,
		Detached:                    detached,
		DisplayTypes:                displayTypes,
		DryRun:                      None[configdomain.DryRun](),
		UnknownBranchType:           unknownBranchType,
		DevRemote:                   devRemote,
		FeatureRegex:                featureRegex,
		ForgeType:                   forgeType,
		GithubConnectorType:         githubConnectorType,
		GithubToken:                 None[forgedomain.GithubToken](),
		GitlabConnectorType:         gitlabConnectorType,
		GitlabToken:                 None[forgedomain.GitlabToken](),
		GitUserEmail:                None[gitdomain.GitUserEmail](),
		GitUserName:                 None[gitdomain.GitUserName](),
		GiteaToken:                  None[forgedomain.GiteaToken](),
		HostingOriginHostname:       hostingOriginHostname,
		Lineage:                     configdomain.NewLineage(),
		MainBranch:                  mainBranch,
		NewBranchType:               newBranchType,
		AutoResolve:                 autoResolve,
		ObservedRegex:               observedRegex,
		Offline:                     None[configdomain.Offline](),
		Order:                       order,
		PerennialBranches:           perennialBranches,
		PerennialRegex:              perennialRegex,
		ProposalBreadcrumb:          proposalBreadcrumb,
		ProposalBreadcrumbDirection: proposalBreadcrumbDirection,
		PushBranches:                pushBranches,
		PushHook:                    pushHook,
		ShareNewBranches:            shareNewBranches,
		ShipDeleteTrackingBranch:    shipDeleteTrackingBranch,
		IgnoreUncommitted:           ignoreUncommitted,
		ShipStrategy:                shipStrategy,
		Stash:                       stash,
		SyncFeatureStrategy:         syncFeatureStrategy,
		SyncPerennialStrategy:       syncPerennialStrategy,
		SyncPrototypeStrategy:       syncPrototypeStrategy,
		SyncTags:                    syncTags,
		SyncUpstream:                syncUpstream,
		Verbose:                     None[configdomain.Verbose](),
	}, ec.Err
}
