// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/git-town/git-town/v11/src/config/confighelpers"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type Config struct {
	configdomain.Access                                // access to the Git configuration settings
	configdomain.FullConfig                            // the merged configuration data
	configFile              configdomain.PartialConfig // content of git-town.toml
	GlobalGitConfig         configdomain.PartialConfig // content of the global Git configuration
	LocalGitConfig          configdomain.PartialConfig // content of the local Git configuration
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
	return confighelpers.DetermineOriginURL(text, self.CodeHostingOriginHostname, self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *Config) OriginURLString() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	output, _ := self.Query("git", "remote", "get-url", gitdomain.OriginRemote.String())
	return strings.TrimSpace(output)
}

func (self *Config) Reload() {
	_, self.GlobalGitConfig, _ = self.LoadCache(true) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	_, self.LocalGitConfig, _ = self.LoadCache(false) // we ignore the Git cache here because reloading a config in the middle of a Git Town command doesn't change the cached initial state of the repo
	self.FullConfig = configdomain.DefaultConfig()
	// TODO: merge this code with the similar code in NewGitTown.
	self.FullConfig.Merge(self.configFile)
	self.FullConfig.Merge(self.GlobalGitConfig)
	self.FullConfig.Merge(self.LocalGitConfig)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *Config) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	slice.Remove(&self.FullConfig.PerennialBranches, branch)
	return self.SetPerennialBranches(self.FullConfig.PerennialBranches)
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *Config) RemoveParent(branch gitdomain.LocalBranchName) {
	self.LocalGitConfig.Lineage.RemoveBranch(branch)
	_ = self.RemoveLocalConfigValue(configdomain.NewParentKey(branch))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *Config) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.MainBranch = branch
	self.LocalGitConfig.MainBranch = &branch
	return self.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *Config) SetNewBranchPush(value configdomain.NewBranchPush, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.NewBranchPush = value
	if global {
		self.GlobalGitConfig.NewBranchPush = &value
		return self.SetGlobalConfigValue(configdomain.KeyPushNewBranches, setting)
	}
	self.LocalGitConfig.NewBranchPush = &value
	return self.SetLocalConfigValue(configdomain.KeyPushNewBranches, setting)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *Config) SetOffline(value configdomain.Offline) error {
	self.FullConfig.Offline = value
	return self.SetGlobalConfigValue(configdomain.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *Config) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage[branch] = parentBranch
	return self.SetLocalConfigValue(configdomain.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *Config) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	return self.SetLocalConfigValue(configdomain.KeyPerennialBranches, branches.Join(" "))
}

// SetPushHook updates the configured push-hook strategy.
func (self *Config) SetPushHookGlobally(value configdomain.PushHook) error {
	self.GlobalGitConfig.PushHook = &value
	self.PushHook = value
	return self.SetGlobalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *Config) SetPushHookLocally(value configdomain.PushHook) error {
	self.LocalGitConfig.PushHook = &value
	self.PushHook = value
	return self.SetLocalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetShipDeleteTrackingBranch updates the configured delete-remote-branch strategy.
func (self *Config) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch) error {
	return self.SetLocalConfigValue(configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *Config) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.LocalGitConfig.SyncFeatureStrategy = &value
	self.FullConfig.SyncFeatureStrategy = value
	return self.SetLocalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

func (self *Config) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.GlobalGitConfig.SyncFeatureStrategy = &value
	self.FullConfig.SyncFeatureStrategy = value
	return self.SetGlobalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *Config) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.LocalGitConfig.SyncPerennialStrategy = &strategy
	self.FullConfig.SyncPerennialStrategy = strategy
	return self.SetLocalConfigValue(configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *Config) SetSyncUpstream(value configdomain.SyncUpstream) error {
	return self.SetLocalConfigValue(configdomain.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

// SetTestOrigin sets the origin to be used for testing.
func (self *Config) SetTestOrigin(value string) error {
	return self.SetLocalConfigValue(configdomain.KeyTestingRemoteURL, value)
}

func NewGitTown(globalConfig, localConfig configdomain.PartialConfig, dryRun bool, runner configdomain.Runner) (*Config, error) {
	configFile, err := configfile.Load()
	if err != nil {
		return nil, err
	}
	config := configdomain.DefaultConfig()
	config.Merge(configFile)
	config.Merge(globalConfig)
	config.Merge(localConfig)
	return &Config{
		Access:          configdomain.Access{Runner: runner},
		FullConfig:      config,
		configFile:      configFile,
		GlobalGitConfig: globalConfig,
		LocalGitConfig:  localConfig,
		DryRun:          dryRun,
		originURLCache:  configdomain.OriginURLCache{},
	}, nil
}
