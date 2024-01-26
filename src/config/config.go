// Package config provides functionality to read and write the Git Town configuration.
// Git Town configuration can exist in a number of locations: in local or global Git metadata or in a configuration file.
// Subspackages implement access to specific configuration locations.
package config

import (
	"strconv"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/git-town/git-town/v11/src/config/confighelpers"
	"github.com/git-town/git-town/v11/src/config/envconfig"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type Config struct {
	GitConfig               gitconfig.Access            // access to the Git configuration settings
	configdomain.FullConfig                             // the merged configuration data
	configFile              *configdomain.PartialConfig // content of git-town.toml, nil = no config file exists
	GlobalGitConfig         configdomain.PartialConfig  // content of the global Git configuration
	LocalGitConfig          configdomain.PartialConfig  // content of the local Git configuration
	DryRun                  bool
	originURLCache          configdomain.OriginURLCache
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *Config) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.PerennialBranches, branches...))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *Config) OriginURL() *giturl.Parts {
	text := self.OriginURLString()
	if text == "" {
		return nil
	}
	return confighelpers.DetermineOriginURL(text, self.HostingOriginHostname, self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *Config) OriginURLString() string {
	remoteOverride := envconfig.OriginURLOverride()
	if remoteOverride != "" {
		return remoteOverride
	}
	return self.GitConfig.OriginRemote()
}

func (self *Config) Reload() {
	_, self.GlobalGitConfig, _ = self.GitConfig.LoadGlobal() // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.GitConfig.LoadLocal()   // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	self.FullConfig = configdomain.DefaultConfig()
	// TODO: merge this code with the similar code in NewConfig.
	if self.configFile != nil {
		self.FullConfig.Merge(*self.configFile)
	}
	self.FullConfig.Merge(self.GlobalGitConfig)
	self.FullConfig.Merge(self.LocalGitConfig)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	self.FullConfig.PerennialBranches = slice.Remove(self.FullConfig.PerennialBranches, branch)
	return self.SetPerennialBranches(self.FullConfig.PerennialBranches)
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *Config) RemoveParent(branch gitdomain.LocalBranchName) {
	self.LocalGitConfig.Lineage.RemoveBranch(branch)
	_ = self.GitConfig.RemoveLocalConfigValue(gitconfig.NewParentKey(branch))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.MainBranch = branch
	if self.configFile != nil {
		self.configFile.MainBranch = &branch
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.MainBranch = &branch
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, branch.String())
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *Config) SetNewBranchPush(value configdomain.NewBranchPush, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.NewBranchPush = value
	if global {
		self.GlobalGitConfig.NewBranchPush = &value
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushNewBranches, setting)
	}
	if self.configFile != nil {
		self.configFile.NewBranchPush = &value
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.NewBranchPush = &value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushNewBranches, setting)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *Config) SetOffline(value configdomain.Offline) error {
	self.FullConfig.Offline = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *Config) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage[branch] = parentBranch
	return self.GitConfig.SetLocalConfigValue(gitconfig.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *Config) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	if self.configFile != nil {
		self.configFile.PerennialBranches = &branches
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.PerennialBranches = &branches
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, branches.Join(" "))
}

// SetPushHook updates the configured push-hook strategy.
func (self *Config) SetPushHookGlobally(value configdomain.PushHook) error {
	self.GlobalGitConfig.PushHook = &value
	self.PushHook = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *Config) SetPushHookLocally(value configdomain.PushHook) error {
	self.LocalGitConfig.PushHook = &value
	self.PushHook = value
	if self.configFile != nil {
		self.configFile.PushHook = &value
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.PushHook = &value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetShipDeleteTrackingBranch updates the configured delete-tracking-branch strategy.
func (self *Config) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch, global bool) error {
	self.FullConfig.ShipDeleteTrackingBranch = value
	if global {
		self.GlobalGitConfig.ShipDeleteTrackingBranch = &value
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
	}
	if self.configFile != nil {
		self.configFile.ShipDeleteTrackingBranch = &value
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.ShipDeleteTrackingBranch = &value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncBeforeShip(value configdomain.SyncBeforeShip, global bool) error {
	self.FullConfig.SyncBeforeShip = value
	if global {
		self.GlobalGitConfig.SyncBeforeShip = &value
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
	}
	if self.configFile != nil {
		self.configFile.SyncBeforeShip = &value
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.SyncBeforeShip = &value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncBeforeShip, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.FullConfig.SyncFeatureStrategy = value
	if self.configFile != nil {
		self.configFile.SyncFeatureStrategy = &value
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.SyncFeatureStrategy = &value
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

func (self *Config) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.GlobalGitConfig.SyncFeatureStrategy = &value
	self.FullConfig.SyncFeatureStrategy = value
	return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncFeatureStrategy, value.String())
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *Config) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.FullConfig.SyncPerennialStrategy = strategy
	if self.configFile != nil {
		self.configFile.SyncPerennialStrategy = &strategy
		return configfile.Save(self.configFile)
	}
	self.LocalGitConfig.SyncPerennialStrategy = &strategy
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *Config) SetSyncUpstream(value configdomain.SyncUpstream, global bool) error {
	self.FullConfig.SyncUpstream = value
	if global {
		self.GlobalGitConfig.SyncUpstream = &value
		return self.GitConfig.SetGlobalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	}
	if self.configFile != nil {
		self.configFile.SyncUpstream = &value
		return configfile.Save(self.configFile)
	}
	return self.GitConfig.SetLocalConfigValue(gitconfig.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

func NewConfig(globalConfig, localConfig configdomain.PartialConfig, dryRun bool, runner gitconfig.Runner) (*Config, error) {
	configFile, err := configfile.Load()
	if err != nil {
		return nil, err
	}
	config := configdomain.DefaultConfig()
	if configFile != nil {
		config.Merge(*configFile)
	}
	config.Merge(globalConfig)
	config.Merge(localConfig)
	return &Config{
		GitConfig:       gitconfig.Access{Runner: runner},
		FullConfig:      config,
		configFile:      configFile,
		GlobalGitConfig: globalConfig,
		LocalGitConfig:  localConfig,
		DryRun:          dryRun,
		originURLCache:  configdomain.OriginURLCache{},
	}, nil
}
