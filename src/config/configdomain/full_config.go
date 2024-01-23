package configdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// FullConfig is the merged configuration to be used by Git Town commands.
type FullConfig struct {
	Aliases                   Aliases
	CodeHostingOriginHostname CodeHostingOriginHostname
	CodeHostingPlatformName   CodeHostingPlatformName
	GiteaToken                GiteaToken
	GitHubToken               GitHubToken
	GitLabToken               GitLabToken
	GitUserEmail              string
	GitUserName               string
	Lineage                   Lineage
	MainBranch                gitdomain.LocalBranchName
	NewBranchPush             NewBranchPush
	Offline                   Offline
	PerennialBranches         gitdomain.LocalBranchNames
	PushHook                  PushHook
	ShipDeleteTrackingBranch  ShipDeleteTrackingBranch
	SyncBeforeShip            SyncBeforeShip
	SyncFeatureStrategy       SyncFeatureStrategy
	SyncPerennialStrategy     SyncPerennialStrategy
	SyncUpstream              SyncUpstream
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *FullConfig) ContainsLineage() bool {
	return len(self.Lineage) > 0
}

// HostingService provides the type-safe name of the code hosting connector to use.
// This function caches its result and can be queried repeatedly.
func (self *FullConfig) HostingService() (Hosting, error) {
	return NewHosting(self.CodeHostingPlatformName)
}

func (self *FullConfig) IsFeatureBranch(branch gitdomain.LocalBranchName) bool {
	return !self.IsMainBranch(branch) && !self.IsPerennialBranch(branch)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *FullConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

func (self *FullConfig) IsOnline() bool {
	return self.Online().Bool()
}

func (self *FullConfig) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return slice.Contains(self.PerennialBranches, branch)
}

func (self *FullConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

// Merges the given PartialConfig into this configuration object.
func (self *FullConfig) Merge(other PartialConfig) {
	for key, value := range other.Aliases {
		self.Aliases[key] = value
	}
	if other.Lineage != nil {
		for child, parent := range *other.Lineage {
			self.Lineage[child] = parent
		}
	}
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
	if other.GitUserEmail != nil {
		self.GitUserEmail = *other.GitUserEmail
	}
	if other.GitUserName != nil {
		self.GitUserName = *other.GitUserName
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

func (self *FullConfig) NoPushHook() NoPushHook {
	return self.PushHook.Negate()
}

func (self *FullConfig) Online() Online {
	return self.Offline.ToOnline()
}

func (self *FullConfig) ShouldNewBranchPush() bool {
	return self.NewBranchPush.Bool()
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() FullConfig {
	return FullConfig{
		Aliases:                   Aliases{},
		CodeHostingOriginHostname: "",
		CodeHostingPlatformName:   "",
		GiteaToken:                "",
		GitLabToken:               "",
		GitHubToken:               "",
		GitUserEmail:              "",
		GitUserName:               "",
		Lineage:                   Lineage{},
		MainBranch:                gitdomain.EmptyLocalBranchName(),
		NewBranchPush:             false,
		Offline:                   false,
		PerennialBranches:         gitdomain.NewLocalBranchNames(),
		PushHook:                  true,
		ShipDeleteTrackingBranch:  true,
		SyncBeforeShip:            false,
		SyncFeatureStrategy:       SyncFeatureStrategyMerge,
		SyncPerennialStrategy:     SyncPerennialStrategyRebase,
		SyncUpstream:              true,
	}
}
