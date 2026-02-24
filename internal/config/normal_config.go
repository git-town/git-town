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
	Aliases                     configdomain.Aliases
	AutoResolve                 configdomain.AutoResolve
	AutoSync                    configdomain.AutoSync
	BitbucketAppPassword        Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername           Option[forgedomain.BitbucketUsername]
	BranchPrefix                Option[configdomain.BranchPrefix]
	BranchTypeOverrides         configdomain.BranchTypeOverrides
	Browser                     Option[configdomain.Browser]
	ContributionRegex           Option[configdomain.ContributionRegex]
	Detached                    configdomain.Detached
	DevRemote                   gitdomain.Remote
	DisplayTypes                configdomain.DisplayTypes
	DryRun                      configdomain.DryRun // whether to only print the Git commands but not execute them
	FeatureRegex                Option[configdomain.FeatureRegex]
	ForgeType                   Option[forgedomain.ForgeType] // None = auto-detect
	ForgejoToken                Option[forgedomain.ForgejoToken]
	GitUserEmail                Option[gitdomain.GitUserEmail]
	GitUserName                 Option[gitdomain.GitUserName]
	GiteaToken                  Option[forgedomain.GiteaToken]
	GithubConnectorType         Option[forgedomain.GithubConnectorType]
	GithubToken                 Option[forgedomain.GithubToken]
	GitlabConnectorType         Option[forgedomain.GitlabConnectorType]
	GitlabToken                 Option[forgedomain.GitlabToken]
	HostingOriginHostname       Option[configdomain.HostingOriginHostname]
	IgnoreUncommitted           configdomain.IgnoreUncommitted
	Lineage                     configdomain.Lineage
	NewBranchType               Option[configdomain.NewBranchType]
	ObservedRegex               Option[configdomain.ObservedRegex]
	Offline                     configdomain.Offline
	Order                       configdomain.Order
	PerennialBranches           gitdomain.LocalBranchNames
	PerennialRegex              Option[configdomain.PerennialRegex]
	ProposalBreadcrumb          configdomain.ProposalBreadcrumb
	ProposalBreadcrumbDirection configdomain.ProposalBreadcrumbDirection
	PushBranches                configdomain.PushBranches
	PushHook                    configdomain.PushHook
	ShareNewBranches            configdomain.ShareNewBranches
	ShipDeleteTrackingBranch    configdomain.ShipDeleteTrackingBranch
	ShipStrategy                configdomain.ShipStrategy
	Stash                       configdomain.Stash
	SyncFeatureStrategy         configdomain.SyncFeatureStrategy
	SyncPerennialStrategy       configdomain.SyncPerennialStrategy
	SyncPrototypeStrategy       configdomain.SyncPrototypeStrategy
	SyncTags                    configdomain.SyncTags
	SyncUpstream                configdomain.SyncUpstream
	UnknownBranchType           configdomain.UnknownBranchType
	Verbose                     configdomain.Verbose
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

// OverwriteWith provides a new NormalConfig that contains data from the given PartialConfig,
// backfilled with data from this NormalConfig where missing
func (self *NormalConfig) OverwriteWith(other configdomain.PartialConfig) NormalConfig {
	return NormalConfig{
		Aliases:                     other.Aliases,
		AutoResolve:                 other.AutoResolve.GetOr(self.AutoResolve),
		AutoSync:                    other.AutoSync.GetOr(self.AutoSync),
		BitbucketAppPassword:        other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:           other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchPrefix:                other.BranchPrefix.Or(self.BranchPrefix),
		BranchTypeOverrides:         other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		Browser:                     other.Browser.Or(self.Browser),
		ContributionRegex:           other.ContributionRegex.Or(self.ContributionRegex),
		Detached:                    other.Detached.GetOr(self.Detached),
		DevRemote:                   other.DevRemote.GetOr(self.DevRemote),
		DisplayTypes:                other.DisplayTypes.GetOr(self.DisplayTypes),
		DryRun:                      other.DryRun.GetOr(self.DryRun),
		FeatureRegex:                other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                   other.ForgeType.Or(self.ForgeType),
		ForgejoToken:                other.ForgejoToken.Or(self.ForgejoToken),
		GitUserEmail:                other.GitUserEmail.Or(self.GitUserEmail),
		GitUserName:                 other.GitUserName.Or(self.GitUserName),
		GiteaToken:                  other.GiteaToken.Or(self.GiteaToken),
		GithubConnectorType:         other.GithubConnectorType.Or(self.GithubConnectorType),
		GithubToken:                 other.GithubToken.Or(self.GithubToken),
		GitlabConnectorType:         other.GitlabConnectorType.Or(self.GitlabConnectorType),
		GitlabToken:                 other.GitlabToken.Or(self.GitlabToken),
		HostingOriginHostname:       other.HostingOriginHostname.Or(self.HostingOriginHostname),
		IgnoreUncommitted:           other.IgnoreUncommitted.GetOr(self.IgnoreUncommitted),
		Lineage:                     other.Lineage.Merge(self.Lineage),
		NewBranchType:               other.NewBranchType.Or(self.NewBranchType),
		ObservedRegex:               other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                     other.Offline.GetOr(self.Offline),
		Order:                       other.Order.GetOr(self.Order),
		PerennialBranches:           other.PerennialBranches.AppendAllMissing(self.PerennialBranches),
		PerennialRegex:              other.PerennialRegex.Or(self.PerennialRegex),
		ProposalBreadcrumb:          other.ProposalBreadcrumb.GetOr(self.ProposalBreadcrumb),
		ProposalBreadcrumbDirection: other.ProposalBreadcrumbDirection.GetOr(self.ProposalBreadcrumbDirection),
		PushBranches:                other.PushBranches.GetOr(self.PushBranches),
		PushHook:                    other.PushHook.GetOr(self.PushHook),
		ShareNewBranches:            other.ShareNewBranches.GetOr(self.ShareNewBranches),
		ShipDeleteTrackingBranch:    other.ShipDeleteTrackingBranch.GetOr(self.ShipDeleteTrackingBranch),
		ShipStrategy:                other.ShipStrategy.GetOr(self.ShipStrategy),
		Stash:                       other.Stash.GetOr(self.Stash),
		SyncFeatureStrategy:         other.SyncFeatureStrategy.GetOr(self.SyncFeatureStrategy),
		SyncPerennialStrategy:       other.SyncPerennialStrategy.GetOr(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:       other.SyncPrototypeStrategy.GetOr(self.SyncPrototypeStrategy),
		SyncTags:                    other.SyncTags.GetOr(self.SyncTags),
		SyncUpstream:                other.SyncUpstream.GetOr(self.SyncUpstream),
		UnknownBranchType:           other.UnknownBranchType.GetOr(self.UnknownBranchType),
		Verbose:                     other.Verbose.GetOr(self.Verbose),
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
			finalMessages.Addf(messages.PerennialBranchRemovedParentEntry, perennialBranch)
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
		BranchPrefix:         None[configdomain.BranchPrefix](),
		BranchTypeOverrides:  configdomain.BranchTypeOverrides{},
		Browser:              None[configdomain.Browser](),
		ContributionRegex:    None[configdomain.ContributionRegex](),
		Detached:             false,
		DevRemote:            gitdomain.RemoteOrigin,
		DisplayTypes: configdomain.DisplayTypes{
			Quantifier:  configdomain.QuantifierNo,
			BranchTypes: []configdomain.BranchType{configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch},
		},
		DryRun:                      false,
		FeatureRegex:                None[configdomain.FeatureRegex](),
		ForgeType:                   None[forgedomain.ForgeType](),
		ForgejoToken:                None[forgedomain.ForgejoToken](),
		GitUserEmail:                None[gitdomain.GitUserEmail](),
		GitUserName:                 None[gitdomain.GitUserName](),
		GiteaToken:                  None[forgedomain.GiteaToken](),
		GithubConnectorType:         None[forgedomain.GithubConnectorType](),
		GithubToken:                 None[forgedomain.GithubToken](),
		GitlabConnectorType:         None[forgedomain.GitlabConnectorType](),
		GitlabToken:                 None[forgedomain.GitlabToken](),
		HostingOriginHostname:       None[configdomain.HostingOriginHostname](),
		IgnoreUncommitted:           false,
		Lineage:                     configdomain.NewLineage(),
		NewBranchType:               None[configdomain.NewBranchType](),
		ObservedRegex:               None[configdomain.ObservedRegex](),
		Offline:                     false,
		Order:                       configdomain.OrderAsc,
		PerennialBranches:           gitdomain.LocalBranchNames{},
		PerennialRegex:              None[configdomain.PerennialRegex](),
		ProposalBreadcrumb:          configdomain.ProposalBreadcrumbNone,
		ProposalBreadcrumbDirection: configdomain.ProposalBreadcrumbDirectionDown,
		PushBranches:                true,
		PushHook:                    true,
		ShareNewBranches:            configdomain.ShareNewBranchesNone,
		ShipDeleteTrackingBranch:    true,
		ShipStrategy:                configdomain.ShipStrategyAPI,
		Stash:                       true,
		SyncFeatureStrategy:         configdomain.SyncFeatureStrategyMerge,
		SyncPerennialStrategy:       configdomain.SyncPerennialStrategyRebase,
		SyncPrototypeStrategy:       configdomain.SyncPrototypeStrategyRebase,
		SyncTags:                    true,
		SyncUpstream:                true,
		UnknownBranchType:           configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
		Verbose:                     false,
	}
}

func NewNormalConfigFromPartial(partial configdomain.PartialConfig, defaults NormalConfig) NormalConfig {
	syncFeatureStrategy := partial.SyncFeatureStrategy.GetOr(defaults.SyncFeatureStrategy)
	proposalBreadcrumbDirection := partial.ProposalBreadcrumbDirection.GetOr(defaults.ProposalBreadcrumbDirection)
	return NormalConfig{
		Aliases:                     partial.Aliases,
		AutoResolve:                 partial.AutoResolve.GetOr(defaults.AutoResolve),
		AutoSync:                    partial.AutoSync.GetOr(defaults.AutoSync),
		BitbucketAppPassword:        partial.BitbucketAppPassword,
		BitbucketUsername:           partial.BitbucketUsername,
		BranchPrefix:                partial.BranchPrefix,
		BranchTypeOverrides:         partial.BranchTypeOverrides,
		Browser:                     partial.Browser.Or(defaults.Browser),
		ContributionRegex:           partial.ContributionRegex,
		Detached:                    partial.Detached.GetOr(defaults.Detached),
		DevRemote:                   partial.DevRemote.GetOr(defaults.DevRemote),
		DisplayTypes:                partial.DisplayTypes.GetOr(defaults.DisplayTypes),
		DryRun:                      partial.DryRun.GetOr(defaults.DryRun),
		FeatureRegex:                partial.FeatureRegex,
		ForgeType:                   partial.ForgeType,
		ForgejoToken:                partial.ForgejoToken,
		GitUserEmail:                partial.GitUserEmail,
		GitUserName:                 partial.GitUserName,
		GiteaToken:                  partial.GiteaToken,
		GithubConnectorType:         partial.GithubConnectorType,
		GithubToken:                 partial.GithubToken,
		GitlabConnectorType:         partial.GitlabConnectorType,
		GitlabToken:                 partial.GitlabToken,
		HostingOriginHostname:       partial.HostingOriginHostname,
		IgnoreUncommitted:           partial.IgnoreUncommitted.GetOr(defaults.IgnoreUncommitted),
		Lineage:                     partial.Lineage,
		NewBranchType:               partial.NewBranchType.Or(defaults.NewBranchType),
		ObservedRegex:               partial.ObservedRegex,
		Offline:                     partial.Offline.GetOr(defaults.Offline),
		Order:                       partial.Order.GetOr(defaults.Order),
		PerennialBranches:           partial.PerennialBranches,
		PerennialRegex:              partial.PerennialRegex,
		ProposalBreadcrumb:          partial.ProposalBreadcrumb.GetOr(defaults.ProposalBreadcrumb),
		ProposalBreadcrumbDirection: proposalBreadcrumbDirection,
		PushBranches:                partial.PushBranches.GetOr(defaults.PushBranches),
		PushHook:                    partial.PushHook.GetOr(defaults.PushHook),
		ShareNewBranches:            partial.ShareNewBranches.GetOr(defaults.ShareNewBranches),
		ShipDeleteTrackingBranch:    partial.ShipDeleteTrackingBranch.GetOr(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:                partial.ShipStrategy.GetOr(defaults.ShipStrategy),
		Stash:                       partial.Stash.GetOr(defaults.Stash),
		SyncFeatureStrategy:         syncFeatureStrategy,
		SyncPerennialStrategy:       partial.SyncPerennialStrategy.GetOr(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:       partial.SyncPrototypeStrategy.GetOr(configdomain.NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                    partial.SyncTags.GetOr(defaults.SyncTags),
		SyncUpstream:                partial.SyncUpstream.GetOr(defaults.SyncUpstream),
		UnknownBranchType:           partial.UnknownBranchType.GetOr(configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)),
		Verbose:                     partial.Verbose.GetOr(defaults.Verbose),
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
