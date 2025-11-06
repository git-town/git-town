package config

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/envconfig"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

// NormalConfig contains the final configuration data to be used by Git Town,
// merged from the various configuration sources:
// - local Git metadata,
// - global Git metadata,
// - configuration file
// - CLI arguments
// - default values
//
// It only contains the configuration data that doesn't need to be prompted from the user if missing,
// but can be used as-is.
// Configuration data that needs to be prompted from the user exists in UnvalidatedConfigData/ValidatedConfigData.
type NormalConfig struct {
	Aliases                  configdomain.Aliases
	AutoResolve              configdomain.AutoResolve
	AutoSync                 configdomain.AutoSync
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      configdomain.BranchTypeOverrides
	ContributionRegex        Option[configdomain.ContributionRegex]
	Detached                 configdomain.Detached
	DevRemote                gitdomain.Remote
	DisplayTypes             configdomain.DisplayTypes
	DryRun                   configdomain.DryRun // whether to only print the Git commands but not execute them
	FeatureRegex             Option[configdomain.FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType] // None = auto-detect
	ForgejoToken             Option[forgedomain.ForgejoToken]
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GitUserEmail             Option[gitdomain.GitUserEmail]
	GitUserName              Option[gitdomain.GitUserName]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[configdomain.HostingOriginHostname]
	Lineage                  configdomain.Lineage
	NewBranchType            Option[configdomain.NewBranchType]
	ObservedRegex            Option[configdomain.ObservedRegex]
	Offline                  configdomain.Offline
	Order                    configdomain.Order
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[configdomain.PerennialRegex]
	ProposalsShowLineage     forgedomain.ProposalsShowLineage
	ProposeBodyTemplate      Option[gitdomain.ProposalBodyTemplate]
	ProposeBodyTemplateFile  Option[gitdomain.ProposalBodyTemplateFile]
	ProposeTitle             Option[configdomain.ProposeTitle]
	PushBranches             configdomain.PushBranches
	PushHook                 configdomain.PushHook
	ShareNewBranches         configdomain.ShareNewBranches
	ShipDeleteTrackingBranch configdomain.ShipDeleteTrackingBranch
	ShipStrategy             configdomain.ShipStrategy
	Stash                    configdomain.Stash
	SyncFeatureStrategy      configdomain.SyncFeatureStrategy
	SyncPerennialStrategy    configdomain.SyncPerennialStrategy
	SyncPrototypeStrategy    configdomain.SyncPrototypeStrategy
	SyncTags                 configdomain.SyncTags
	SyncUpstream             configdomain.SyncUpstream
	UnknownBranchType        configdomain.UnknownBranchType
	Verbose                  configdomain.Verbose
}

// Author provides the locally Git configured user.
func (self *NormalConfig) Author() Option[gitdomain.Author] {
	email, hasEmail := self.GitUserEmail.Get()
	name, hasName := self.GitUserName.Get()
	if hasEmail && hasName {
		return Some(gitdomain.Author(fmt.Sprintf("%s <%s>", name, email)))
	}
	return None[gitdomain.Author]()
}

// DevURL provides the URL for the development remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) DevURL(querier subshelldomain.Querier) Option[giturl.Parts] {
	return self.RemoteURL(querier, self.DevRemote)
}

// provides a new NormalConfig that contains data from the given PartialConfig,
// backfilled with data from this NormalConfig where missing
func (self *NormalConfig) OverwriteWith(other configdomain.PartialConfig) NormalConfig {
	return NormalConfig{
		Aliases:                  other.Aliases,
		AutoResolve:              other.AutoResolve.GetOr(self.AutoResolve),
		AutoSync:                 other.AutoSync.GetOr(self.AutoSync),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		Detached:                 other.Detached.GetOr(self.Detached),
		DevRemote:                other.DevRemote.GetOr(self.DevRemote),
		DisplayTypes:             other.DisplayTypes.GetOr(self.DisplayTypes),
		DryRun:                   other.DryRun.GetOr(self.DryRun),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                other.ForgeType.Or(self.ForgeType),
		ForgejoToken:             other.ForgejoToken.Or(self.ForgejoToken),
		GitHubConnectorType:      other.GitHubConnectorType.Or(self.GitHubConnectorType),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitLabConnectorType:      other.GitLabConnectorType.Or(self.GitLabConnectorType),
		GitLabToken:              other.GitLabToken.Or(self.GitLabToken),
		GitUserEmail:             other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:              other.GitUserName.Or(self.GitUserName),
		GiteaToken:               other.GiteaToken.Or(self.GiteaToken),
		HostingOriginHostname:    other.HostingOriginHostname.Or(self.HostingOriginHostname),
		Lineage:                  other.Lineage.Merge(self.Lineage),
		NewBranchType:            other.NewBranchType.Or(self.NewBranchType),
		ObservedRegex:            other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                  other.Offline.GetOr(self.Offline),
		Order:                    other.Order.GetOr(self.Order),
		PerennialBranches:        other.PerennialBranches.AppendAllMissing(self.PerennialBranches),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		ProposalsShowLineage:     other.ProposalsShowLineage.GetOr(self.ProposalsShowLineage),
		ProposeBodyTemplate:      other.ProposeBodyTemplate.Or(self.ProposeBodyTemplate),
		ProposeBodyTemplateFile:  other.ProposeBodyTemplateFile.Or(self.ProposeBodyTemplateFile),
		ProposeTitle:             other.ProposeTitle.Or(self.ProposeTitle),
		PushBranches:             other.PushBranches.GetOr(self.PushBranches),
		PushHook:                 other.PushHook.GetOr(self.PushHook),
		ShareNewBranches:         other.ShareNewBranches.GetOr(self.ShareNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.GetOr(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.GetOr(self.ShipStrategy),
		Stash:                    other.Stash.GetOr(self.Stash),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.GetOr(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.GetOr(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.GetOr(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.GetOr(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.GetOr(self.SyncUpstream),
		UnknownBranchType:        other.UnknownBranchType.GetOr(self.UnknownBranchType),
		Verbose:                  other.Verbose.GetOr(self.Verbose),
	}
}

func (self *NormalConfig) PartialBranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	// check the branch type overrides
	if branchTypeOverride, hasBranchTypeOverride := self.BranchTypeOverrides[branch]; hasBranchTypeOverride {
		return branchTypeOverride
	}
	// check the configured branch lists
	if slices.Contains(self.PerennialBranches, branch) {
		return configdomain.BranchTypePerennialBranch
	}
	// check if a regex matches
	if regex, has := self.ContributionRegex.Get(); has && regex.MatchesBranch(branch) {
		return configdomain.BranchTypeContributionBranch
	}
	if regex, has := self.FeatureRegex.Get(); has && regex.MatchesBranch(branch) {
		return configdomain.BranchTypeFeatureBranch
	}
	if regex, has := self.ObservedRegex.Get(); has && regex.MatchesBranch(branch) {
		return configdomain.BranchTypeObservedBranch
	}
	if regex, has := self.PerennialRegex.Get(); has && regex.MatchesBranch(branch) {
		return configdomain.BranchTypePerennialBranch
	}
	// branch doesn't match any of the overrides --> unknown branch type
	return self.UnknownBranchType.BranchType()
}

func (self *NormalConfig) PartialBranchesOfType(branchType configdomain.BranchType) gitdomain.LocalBranchNames {
	matching := set.New[gitdomain.LocalBranchName]()
	switch branchType {
	case configdomain.BranchTypeContributionBranch:
	case configdomain.BranchTypeFeatureBranch:
	case configdomain.BranchTypeMainBranch:
		// main branch is stored in ValidatedConfig
	case configdomain.BranchTypeObservedBranch:
	case configdomain.BranchTypeParkedBranch:
	case configdomain.BranchTypePerennialBranch:
		matching.Add(self.PerennialBranches...)
	case configdomain.BranchTypePrototypeBranch:
	}
	for key, value := range self.BranchTypeOverrides { // okay to iterate the map in random order here because we add to a set
		if value == branchType {
			matching.Add(key)
		}
	}
	return matching.Values()
}

// RemoteURL provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) RemoteURL(querier subshelldomain.Querier, remote gitdomain.Remote) Option[giturl.Parts] {
	urlStr, hasURLStr := remoteURLString(querier, remote).Get()
	if !hasURLStr {
		return None[giturl.Parts]()
	}
	url, hasURL := giturl.Parse(urlStr).Get()
	if !hasURL {
		return None[giturl.Parts]()
	}
	if hostnameOverride, hasHostNameOverride := self.HostingOriginHostname.Get(); hasHostNameOverride {
		url.Host = hostnameOverride.String()
	}
	return Some(url)
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveParent(runner, branch)
}

func (self *NormalConfig) RemovePerennialAncestors(runner subshelldomain.Runner, finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = gitconfig.RemoveParent(runner, perennialBranch)
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(runner subshelldomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return gitconfig.SetParent(runner, branch, parentBranch)
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *NormalConfig) SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	var err error
	if slices.Compare(self.PerennialBranches, branches) != 0 {
		err = gitconfig.SetPerennialBranches(runner, branches, configdomain.ConfigScopeLocal)
	}
	self.PerennialBranches = branches
	return err
}

func DefaultNormalConfig() NormalConfig {
	return NormalConfig{
		Aliases:              configdomain.Aliases{},
		AutoResolve:          true,
		AutoSync:             true,
		BitbucketAppPassword: None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:    None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:  configdomain.BranchTypeOverrides{},
		ContributionRegex:    None[configdomain.ContributionRegex](),
		Detached:             false,
		DevRemote:            gitdomain.RemoteOrigin,
		DisplayTypes: configdomain.DisplayTypes{
			Quantifier:  configdomain.QuantifierNo,
			BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch},
		},
		DryRun:                   false,
		FeatureRegex:             None[configdomain.FeatureRegex](),
		ForgeType:                None[forgedomain.ForgeType](),
		ForgejoToken:             None[forgedomain.ForgejoToken](),
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
		GitLabToken:              None[forgedomain.GitLabToken](),
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
		Lineage:                  configdomain.NewLineage(),
		NewBranchType:            None[configdomain.NewBranchType](),
		ObservedRegex:            None[configdomain.ObservedRegex](),
		Offline:                  false,
		Order:                    configdomain.OrderAsc,
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[configdomain.PerennialRegex](),
		ProposalsShowLineage:     forgedomain.ProposalsShowLineageNone,
		ProposeBodyTemplate:      None[gitdomain.ProposalBodyTemplate](),
		ProposeBodyTemplateFile:  None[gitdomain.ProposalBodyTemplateFile](),
		ProposeTitle:             None[configdomain.ProposeTitle](),
		PushBranches:             true,
		PushHook:                 true,
		ShareNewBranches:         configdomain.ShareNewBranchesNone,
		ShipDeleteTrackingBranch: true,
		ShipStrategy:             configdomain.ShipStrategyAPI,
		Stash:                    true,
		SyncFeatureStrategy:      configdomain.SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    configdomain.SyncPerennialStrategyRebase,
		SyncPrototypeStrategy:    configdomain.SyncPrototypeStrategyRebase,
		SyncTags:                 true,
		SyncUpstream:             true,
		UnknownBranchType:        configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
		Verbose:                  false,
	}
}

func NewNormalConfigFromPartial(partial configdomain.PartialConfig, defaults NormalConfig) NormalConfig {
	syncFeatureStrategy := partial.SyncFeatureStrategy.GetOr(defaults.SyncFeatureStrategy)
	return NormalConfig{
		Aliases:                  partial.Aliases,
		AutoResolve:              partial.AutoResolve.GetOr(defaults.AutoResolve),
		AutoSync:                 partial.AutoSync.GetOr(defaults.AutoSync),
		BitbucketAppPassword:     partial.BitbucketAppPassword,
		BitbucketUsername:        partial.BitbucketUsername,
		BranchTypeOverrides:      partial.BranchTypeOverrides,
		ContributionRegex:        partial.ContributionRegex,
		Detached:                 partial.Detached.GetOr(defaults.Detached),
		DevRemote:                partial.DevRemote.GetOr(defaults.DevRemote),
		DisplayTypes:             partial.DisplayTypes.GetOr(defaults.DisplayTypes),
		DryRun:                   partial.DryRun.GetOr(defaults.DryRun),
		FeatureRegex:             partial.FeatureRegex,
		ForgeType:                partial.ForgeType,
		ForgejoToken:             partial.ForgejoToken,
		GitHubConnectorType:      partial.GitHubConnectorType,
		GitHubToken:              partial.GitHubToken,
		GitLabConnectorType:      partial.GitLabConnectorType,
		GitLabToken:              partial.GitLabToken,
		GitUserEmail:             partial.GitUserEmail,
		GitUserName:              partial.GitUserName,
		GiteaToken:               partial.GiteaToken,
		HostingOriginHostname:    partial.HostingOriginHostname,
		Lineage:                  partial.Lineage,
		NewBranchType:            partial.NewBranchType.Or(defaults.NewBranchType),
		ObservedRegex:            partial.ObservedRegex,
		Offline:                  partial.Offline.GetOr(defaults.Offline),
		Order:                    partial.Order.GetOr(defaults.Order),
		PerennialBranches:        partial.PerennialBranches,
		PerennialRegex:           partial.PerennialRegex,
		ProposalsShowLineage:     partial.ProposalsShowLineage.GetOr(defaults.ProposalsShowLineage),
		ProposeBodyTemplate:      partial.ProposeBodyTemplate.Or(defaults.ProposeBodyTemplate),
		ProposeBodyTemplateFile:  partial.ProposeBodyTemplateFile.Or(defaults.ProposeBodyTemplateFile),
		ProposeTitle:             partial.ProposeTitle.Or(defaults.ProposeTitle),
		PushBranches:             partial.PushBranches.GetOr(defaults.PushBranches),
		PushHook:                 partial.PushHook.GetOr(defaults.PushHook),
		ShareNewBranches:         partial.ShareNewBranches.GetOr(defaults.ShareNewBranches),
		ShipDeleteTrackingBranch: partial.ShipDeleteTrackingBranch.GetOr(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:             partial.ShipStrategy.GetOr(defaults.ShipStrategy),
		Stash:                    partial.Stash.GetOr(defaults.Stash),
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    partial.SyncPerennialStrategy.GetOr(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    partial.SyncPrototypeStrategy.GetOr(configdomain.NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 partial.SyncTags.GetOr(defaults.SyncTags),
		SyncUpstream:             partial.SyncUpstream.GetOr(defaults.SyncUpstream),
		UnknownBranchType:        partial.UnknownBranchType.GetOr(configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)),
		Verbose:                  partial.Verbose.GetOr(defaults.Verbose),
	}
}

// remoteURLString provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func remoteURLString(querier subshelldomain.Querier, remote gitdomain.Remote) Option[string] {
	remoteOverride := envconfig.RemoteURLOverride()
	if remoteOverride.IsSome() {
		return remoteOverride
	}
	return gitconfig.RemoteURL(querier, remote)
}
