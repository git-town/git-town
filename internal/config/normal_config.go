package config

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
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
	BitbucketAppPassword     Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername        Option[forgedomain.BitbucketUsername]
	BranchTypeOverrides      configdomain.BranchTypeOverrides
	CodebergToken            Option[forgedomain.CodebergToken]
	ContributionRegex        Option[configdomain.ContributionRegex]
	Detached                 configdomain.Detached
	DevRemote                gitdomain.Remote
	DryRun                   configdomain.DryRun // whether to only print the Git commands but not execute them
	FeatureRegex             Option[configdomain.FeatureRegex]
	ForgeType                Option[forgedomain.ForgeType] // None = auto-detect
	GitHubConnectorType      Option[forgedomain.GitHubConnectorType]
	GitHubToken              Option[forgedomain.GitHubToken]
	GitLabConnectorType      Option[forgedomain.GitLabConnectorType]
	GitLabToken              Option[forgedomain.GitLabToken]
	GiteaToken               Option[forgedomain.GiteaToken]
	HostingOriginHostname    Option[configdomain.HostingOriginHostname]
	Lineage                  configdomain.Lineage
	NewBranchType            Option[configdomain.NewBranchType]
	ObservedRegex            Option[configdomain.ObservedRegex]
	Offline                  configdomain.Offline
	PerennialBranches        gitdomain.LocalBranchNames
	PerennialRegex           Option[configdomain.PerennialRegex]
	ProposalsShowLineage     forgedomain.ProposalsShowLineage
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

// DevURL provides the URL for the development remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) DevURL(querier subshelldomain.Querier) Option[giturl.Parts] {
	return self.RemoteURL(querier, self.DevRemote)
}

func (self *NormalConfig) NoPushHook() configdomain.NoPushHook {
	return self.PushHook.Negate()
}

// provides a new NormalConfig that contains data from the given PartialConfig,
// backfilled with data from this NormalConfig where missing
func (self *NormalConfig) OverwriteWith(other configdomain.PartialConfig) NormalConfig {
	return NormalConfig{
		Aliases:                  other.Aliases,
		AutoResolve:              other.AutoResolve.GetOrElse(self.AutoResolve),
		BitbucketAppPassword:     other.BitbucketAppPassword.Or(self.BitbucketAppPassword),
		BitbucketUsername:        other.BitbucketUsername.Or(self.BitbucketUsername),
		BranchTypeOverrides:      other.BranchTypeOverrides.Concat(self.BranchTypeOverrides),
		CodebergToken:            other.CodebergToken.Or(self.CodebergToken),
		ContributionRegex:        other.ContributionRegex.Or(self.ContributionRegex),
		Detached:                 other.Detached.GetOrElse(self.Detached),
		DevRemote:                other.DevRemote.GetOrElse(self.DevRemote),
		DryRun:                   other.DryRun.GetOrElse(self.DryRun),
		FeatureRegex:             other.FeatureRegex.Or(self.FeatureRegex),
		ForgeType:                other.ForgeType.Or(self.ForgeType),
		GitHubConnectorType:      other.GitHubConnectorType.Or(self.GitHubConnectorType),
		GitHubToken:              other.GitHubToken.Or(self.GitHubToken),
		GitLabConnectorType:      other.GitLabConnectorType.Or(self.GitLabConnectorType),
		GitLabToken:              other.GitLabToken.Or(self.GitLabToken),
		GiteaToken:               other.GiteaToken.Or(self.GiteaToken),
		HostingOriginHostname:    other.HostingOriginHostname.Or(self.HostingOriginHostname),
		Lineage:                  other.Lineage.Merge(self.Lineage),
		NewBranchType:            other.NewBranchType.Or(self.NewBranchType),
		ObservedRegex:            other.ObservedRegex.Or(self.ObservedRegex),
		Offline:                  other.Offline.GetOrElse(self.Offline),
		PerennialBranches:        other.PerennialBranches.AppendAllMissing(self.PerennialBranches),
		PerennialRegex:           other.PerennialRegex.Or(self.PerennialRegex),
		ProposalsShowLineage:     other.ProposalsShowLineage.GetOrElse(self.ProposalsShowLineage),
		PushHook:                 other.PushHook.GetOrElse(self.PushHook),
		ShareNewBranches:         other.ShareNewBranches.GetOrElse(self.ShareNewBranches),
		ShipDeleteTrackingBranch: other.ShipDeleteTrackingBranch.GetOrElse(self.ShipDeleteTrackingBranch),
		ShipStrategy:             other.ShipStrategy.GetOrElse(self.ShipStrategy),
		Stash:                    other.Stash.GetOrElse(self.Stash),
		SyncFeatureStrategy:      other.SyncFeatureStrategy.GetOrElse(self.SyncFeatureStrategy),
		SyncPerennialStrategy:    other.SyncPerennialStrategy.GetOrElse(self.SyncPerennialStrategy),
		SyncPrototypeStrategy:    other.SyncPrototypeStrategy.GetOrElse(self.SyncPrototypeStrategy),
		SyncTags:                 other.SyncTags.GetOrElse(self.SyncTags),
		SyncUpstream:             other.SyncUpstream.GetOrElse(self.SyncUpstream),
		UnknownBranchType:        other.UnknownBranchType.GetOrElse(self.UnknownBranchType),
		Verbose:                  other.Verbose.GetOrElse(self.Verbose),
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
	for key, value := range self.BranchTypeOverrides {
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
		Aliases:                  configdomain.Aliases{},
		AutoResolve:              true,
		BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
		BitbucketUsername:        None[forgedomain.BitbucketUsername](),
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
		CodebergToken:            None[forgedomain.CodebergToken](),
		ContributionRegex:        None[configdomain.ContributionRegex](),
		Detached:                 false,
		DevRemote:                gitdomain.RemoteOrigin,
		DryRun:                   false,
		FeatureRegex:             None[configdomain.FeatureRegex](),
		ForgeType:                None[forgedomain.ForgeType](),
		GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
		GitHubToken:              None[forgedomain.GitHubToken](),
		GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
		GitLabToken:              None[forgedomain.GitLabToken](),
		GiteaToken:               None[forgedomain.GiteaToken](),
		HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
		Lineage:                  configdomain.NewLineage(),
		NewBranchType:            None[configdomain.NewBranchType](),
		ObservedRegex:            None[configdomain.ObservedRegex](),
		Offline:                  false,
		PerennialBranches:        gitdomain.LocalBranchNames{},
		PerennialRegex:           None[configdomain.PerennialRegex](),
		ProposalsShowLineage:     forgedomain.ProposalsShowLineageNone,
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
	syncFeatureStrategy := partial.SyncFeatureStrategy.GetOrElse(defaults.SyncFeatureStrategy)
	return NormalConfig{
		Aliases:                  partial.Aliases,
		AutoResolve:              partial.AutoResolve.GetOrElse(defaults.AutoResolve),
		BitbucketAppPassword:     partial.BitbucketAppPassword,
		BitbucketUsername:        partial.BitbucketUsername,
		BranchTypeOverrides:      partial.BranchTypeOverrides,
		CodebergToken:            partial.CodebergToken,
		ContributionRegex:        partial.ContributionRegex,
		Detached:                 partial.Detached.GetOrElse(defaults.Detached),
		DevRemote:                partial.DevRemote.GetOrElse(defaults.DevRemote),
		DryRun:                   partial.DryRun.GetOrDefault(),
		FeatureRegex:             partial.FeatureRegex,
		ForgeType:                partial.ForgeType,
		GitHubConnectorType:      partial.GitHubConnectorType,
		GitHubToken:              partial.GitHubToken,
		GitLabConnectorType:      partial.GitLabConnectorType,
		GitLabToken:              partial.GitLabToken,
		GiteaToken:               partial.GiteaToken,
		HostingOriginHostname:    partial.HostingOriginHostname,
		Lineage:                  partial.Lineage,
		NewBranchType:            partial.NewBranchType.Or(defaults.NewBranchType),
		ObservedRegex:            partial.ObservedRegex,
		Offline:                  partial.Offline.GetOrElse(defaults.Offline),
		PerennialBranches:        partial.PerennialBranches,
		PerennialRegex:           partial.PerennialRegex,
		ProposalsShowLineage:     partial.ProposalsShowLineage.GetOrElse(defaults.ProposalsShowLineage),
		PushHook:                 partial.PushHook.GetOrElse(defaults.PushHook),
		ShareNewBranches:         partial.ShareNewBranches.GetOrElse(defaults.ShareNewBranches),
		ShipDeleteTrackingBranch: partial.ShipDeleteTrackingBranch.GetOrElse(defaults.ShipDeleteTrackingBranch),
		ShipStrategy:             partial.ShipStrategy.GetOrElse(defaults.ShipStrategy),
		Stash:                    partial.Stash.GetOrElse(defaults.Stash),
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    partial.SyncPerennialStrategy.GetOrElse(defaults.SyncPerennialStrategy),
		SyncPrototypeStrategy:    partial.SyncPrototypeStrategy.GetOrElse(configdomain.NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy)),
		SyncTags:                 partial.SyncTags.GetOrElse(defaults.SyncTags),
		SyncUpstream:             partial.SyncUpstream.GetOrElse(defaults.SyncUpstream),
		UnknownBranchType:        partial.UnknownBranchType.GetOrElse(configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch)),
		Verbose:                  partial.Verbose.GetOrDefault(),
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
