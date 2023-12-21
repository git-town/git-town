package configdomain

import (
	"github.com/git-town/git-town/v11/src/domain"
)

// Config is the merged configuration to be used by Git Town commands.
type Config struct {
	CodeHostingOriginHostname CodeHostingOriginHostname
	CodeHostingPlatformName   CodeHostingPlatformName
	GiteaToken                GiteaToken
	GitHubToken               GitHubToken
	GitLabToken               GitLabToken
	MainBranch                domain.LocalBranchName
	NewBranchPush             NewBranchPush
	Offline                   Offline
	PerennialBranches         domain.LocalBranchNames
	PushHook                  PushHook
	ShipDeleteTrackingBranch  ShipDeleteTrackingBranch
	SyncBeforeShip            SyncBeforeShip
	SyncFeatureStrategy       SyncFeatureStrategy
	SyncPerennialStrategy     SyncPerennialStrategy
	SyncUpstream              SyncUpstream
}

// Merges the given PartialConfig into this configuration object.
func (self *Config) Merge(other PartialConfig) {
	if other.CodeHostingOriginHostname != nil {
		self.CodeHostingOriginHostname = *other.CodeHostingOriginHostname
	}
	if other.CodeHostingPlatformName != nil {
		self.CodeHostingPlatformName = *other.CodeHostingPlatformName
	}
	if other.GiteaToken != nil {
		self.GiteaToken = *other.GiteaToken
	}
	if other.GitHubToken != nil {
		self.GitHubToken = *other.GitHubToken
	}
	if other.GitLabToken != nil {
		self.GitLabToken = *other.GitLabToken
	}
	if other.MainBranch != nil {
		self.MainBranch = *other.MainBranch
	}
	if other.NewBranchPush != nil {
		self.NewBranchPush = *other.NewBranchPush
	}
	if other.Offline != nil {
		self.Offline = *other.Offline
	}
	if other.PerennialBranches != nil {
		self.PerennialBranches = append(self.PerennialBranches, *other.PerennialBranches...)
	}
	if other.PushHook != nil {
		self.PushHook = *other.PushHook
	}
	if other.ShipDeleteTrackingBranch != nil {
		self.ShipDeleteTrackingBranch = *other.ShipDeleteTrackingBranch
	}
	if other.SyncBeforeShip != nil {
		self.SyncBeforeShip = *other.SyncBeforeShip
	}
	if other.SyncFeatureStrategy != nil {
		self.SyncFeatureStrategy = *other.SyncFeatureStrategy
	}
	if other.SyncPerennialStrategy != nil {
		self.SyncPerennialStrategy = *other.SyncPerennialStrategy
	}
	if other.SyncUpstream != nil {
		self.SyncUpstream = *other.SyncUpstream
	}
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() Config {
	return Config{
		CodeHostingOriginHostname: "",
		CodeHostingPlatformName:   "",
		GiteaToken:                "",
		GitLabToken:               "",
		GitHubToken:               "",
		MainBranch:                domain.EmptyLocalBranchName(),
		NewBranchPush:             false,
		Offline:                   false,
		PerennialBranches:         domain.NewLocalBranchNames(),
		PushHook:                  true,
		ShipDeleteTrackingBranch:  true,
		SyncBeforeShip:            false,
		SyncFeatureStrategy:       SyncFeatureStrategyMerge,
		SyncPerennialStrategy:     SyncPerennialStrategyRebase,
		SyncUpstream:              true,
	}
}
