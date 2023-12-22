package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/confighelpers"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
)

// GitTown provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type GitTown struct {
	configdomain.CachedAccess      // access to the Git configuration settings
	configdomain.Config            // the merged configuration data
	DryRun                    bool // single source of truth for whether to dry-run Git commands in this repo
	originURLCache            configdomain.OriginURLCache
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *GitTown) AddToPerennialBranches(branches ...gitdomain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.PerennialBranches, branches...))
}

func (self *GitTown) BranchTypes() configdomain.BranchTypes {
	return configdomain.BranchTypes{
		MainBranch:        self.Config.MainBranch,
		PerennialBranches: self.Config.PerennialBranches,
	}
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *GitTown) ContainsLineage() bool {
	for key := range self.LocalCache {
		if strings.HasPrefix(key.String(), "git-town-branch.") {
			return true
		}
	}
	return false
}

// HostingService provides the type-safe name of the code hosting connector to use.
// This function caches its result and can be queried repeatedly.
func (self *GitTown) HostingService() (configdomain.Hosting, error) {
	return configdomain.NewHosting(self.Config.CodeHostingPlatformName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *GitTown) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.Config.MainBranch
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *GitTown) OriginURL() *giturl.Parts {
	text := self.OriginURLString()
	if text == "" {
		return nil
	}
	return confighelpers.DetermineOriginURL(text, self.CodeHostingOriginHostname, self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *GitTown) OriginURLString() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	output, _ := self.Query("git", "remote", "get-url", gitdomain.OriginRemote.String())
	return strings.TrimSpace(output)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *GitTown) RemoveFromPerennialBranches(branch gitdomain.LocalBranchName) error {
	slice.Remove(&self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *GitTown) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.MainBranch = branch
	self.LocalConfig.MainBranch = &branch
	return self.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *GitTown) SetNewBranchPush(value configdomain.NewBranchPush, global bool) error {
	setting := strconv.FormatBool(bool(value))
	self.NewBranchPush = value
	if global {
		self.GlobalConfig.NewBranchPush = &value
		return self.SetGlobalConfigValue(configdomain.KeyPushNewBranches, setting)
	}
	self.LocalConfig.NewBranchPush = &value
	return self.SetLocalConfigValue(configdomain.KeyPushNewBranches, setting)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *GitTown) SetOffline(value configdomain.Offline) error {
	self.Config.Offline = value
	return self.SetGlobalConfigValue(configdomain.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *GitTown) SetParent(branch, parentBranch gitdomain.LocalBranchName) error {
	return self.SetLocalConfigValue(configdomain.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *GitTown) SetPerennialBranches(branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	return self.SetLocalConfigValue(configdomain.KeyPerennialBranches, branches.Join(" "))
}

// SetPushHook updates the configured push-hook strategy.
func (self *GitTown) SetPushHookGlobally(value configdomain.PushHook) error {
	self.GlobalConfig.PushHook = &value
	self.PushHook = value
	return self.SetGlobalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(value.Bool()))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *GitTown) SetPushHookLocally(value configdomain.PushHook) error {
	self.LocalConfig.PushHook = &value
	self.PushHook = value
	return self.SetLocalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetShipDeleteTrackingBranch updates the configured delete-remote-branch strategy.
func (self *GitTown) SetShipDeleteTrackingBranch(value configdomain.ShipDeleteTrackingBranch) error {
	return self.SetLocalConfigValue(configdomain.KeyShipDeleteTrackingBranch, strconv.FormatBool(value.Bool()))
}

func (self *GitTown) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	self.LocalConfig.SyncFeatureStrategy = &value
	self.Config.SyncFeatureStrategy = value
	return self.SetLocalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

func (self *GitTown) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	self.GlobalConfig.SyncFeatureStrategy = &value
	self.Config.SyncFeatureStrategy = value
	return self.SetGlobalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *GitTown) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	self.LocalConfig.SyncPerennialStrategy = &strategy
	self.Config.SyncPerennialStrategy = strategy
	return self.SetLocalConfigValue(configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetSyncUpstream updates the configured sync-upstream strategy.
func (self *GitTown) SetSyncUpstream(value configdomain.SyncUpstream) error {
	return self.SetLocalConfigValue(configdomain.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

// SetTestOrigin sets the origin to be used for testing.
func (self *GitTown) SetTestOrigin(value string) error {
	return self.SetLocalConfigValue(configdomain.KeyTestingRemoteURL, value)
}

func NewGitTown(fullCache configdomain.FullCache, runner configdomain.Runner, dryrun bool) (*GitTown, error) {
	configFile, err := configdomain.LoadConfigFile()
	if err != nil {
		return nil, err
	}
	config := configdomain.DefaultConfig()
	config.Merge(configFile)
	config.Merge(fullCache.GlobalConfig)
	config.Merge(fullCache.LocalConfig)
	return &GitTown{
		CachedAccess:   configdomain.NewCachedAccess(fullCache, runner),
		Config:         config,
		DryRun:         dryrun,
		originURLCache: configdomain.OriginURLCache{},
	}, nil
}
