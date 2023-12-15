package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/confighelpers"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/src/messages"
)

// GitTown provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type GitTown struct {
	gitconfig.CachedAccess      // access to the Git configuration settings
	configdomain.Config         // the merged configuration data
	DryRun                 bool // single source of truth for whether to dry-run Git commands in this repo
	originURLCache         configdomain.OriginURLCache
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *GitTown) AddToPerennialBranches(branches ...domain.LocalBranchName) error {
	self.Config.PerennialBranches = append(self.Config.PerennialBranches, branches...)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

func (self *GitTown) BranchTypes() domain.BranchTypes {
	return domain.BranchTypes{
		MainBranch:        self.Config.MainBranch,
		PerennialBranches: self.Config.PerennialBranches,
	}
}

func NewGitTown(fullCache gitconfig.FullCache, runner gitconfig.Runner, dryrun bool) *GitTown {
	config := configdomain.DefaultConfig()
	config.Merge(fullCache.GlobalConfig)
	config.Merge(fullCache.LocalConfig)
	return &GitTown{
		Config:         config,
		CachedAccess:   gitconfig.NewCachedAccess(fullCache, runner),
		DryRun:         dryrun,
		originURLCache: configdomain.OriginURLCache{},
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

// GitAlias provides the currently set alias for the given Git Town command.
func (self *GitTown) GitAlias(alias configdomain.Alias) string {
	return self.GlobalConfigValue(configdomain.NewAliasKey(alias))
}

// HostingService provides the type-safe name of the code hosting connector to use.
// This function caches its result and can be queried repeatedly.
func (self *GitTown) HostingService() (configdomain.Hosting, error) {
	return configdomain.NewHosting(self.Config.CodeHostingPlatformName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *GitTown) IsMainBranch(branch domain.LocalBranchName) bool {
	return branch == self.Config.MainBranch
}

// Lineage provides the configured ancestry information for this Git repo.
func (self *GitTown) Lineage(deleteEntry func(configdomain.Key) error) configdomain.Lineage {
	lineage := configdomain.Lineage{}
	for _, key := range self.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := domain.NewLocalBranchName(strings.TrimSuffix(strings.TrimPrefix(key.String(), "git-town-branch."), ".parent"))
		parentName := self.LocalConfigValue(key)
		if parentName == "" {
			_ = deleteEntry(key)
			fmt.Printf("\nNOTICE: I have found an empty parent configuration entry for branch %q.\n", child)
			fmt.Println("I have deleted this configuration entry.")
		} else {
			parent := domain.NewLocalBranchName(parentName)
			lineage[child] = parent
		}
	}
	return lineage
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (self *GitTown) OriginOverride() configdomain.OriginHostnameOverride {
	return configdomain.OriginHostnameOverride(self.LocalConfigValue(configdomain.KeyCodeHostingOriginHostname))
}

// OriginURL provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *GitTown) OriginURL() *giturl.Parts {
	text := self.OriginURLString()
	if text == "" {
		return nil
	}
	return confighelpers.DetermineOriginURL(text, self.OriginOverride(), self.originURLCache)
}

// OriginURLString provides the URL for the "origin" remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *GitTown) OriginURLString() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	output, _ := self.Query("git", "remote", "get-url", domain.OriginRemote.String())
	return strings.TrimSpace(output)
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHook() (configdomain.PushHook, error) {
	setting := self.LocalOrGlobalConfigValue(configdomain.KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := gohacks.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, configdomain.KeyPushHook, setting)
	}
	return configdomain.PushHook(result), nil
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHookGlobal() (configdomain.PushHook, error) {
	setting := self.GlobalConfigValue(configdomain.KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := gohacks.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueGlobalInvalid, configdomain.KeyPushHook, setting)
	}
	return configdomain.PushHook(result), nil
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *GitTown) RemoveFromPerennialBranches(branch domain.LocalBranchName) error {
	slice.Remove(&self.Config.PerennialBranches, branch)
	return self.SetPerennialBranches(self.Config.PerennialBranches)
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (self *GitTown) RemoveLocalGitConfiguration() error {
	err := self.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 128 {
				// Git returns exit code 128 when trying to delete a non-existing config section.
				// This is not an error condition in this workflow so we can ignore it here.
				return nil
			}
		}
		return fmt.Errorf(messages.ConfigRemoveError, err)
	}
	for _, key := range self.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		err = self.Run("git", "config", "--unset", key.String())
		if err != nil {
			return fmt.Errorf(messages.ConfigRemoveError, err)
		}
	}
	return nil
}

// RemoveMainBranchConfiguration removes the configuration entry for the main branch name.
func (self *GitTown) RemoveMainBranchConfiguration() error {
	return self.RemoveLocalConfigValue(configdomain.KeyMainBranch)
}

// RemoveParent removes the parent branch entry for the given branch
// from the Git configuration.
func (self *GitTown) RemoveParent(branch domain.LocalBranchName) {
	// ignoring errors here because the entry might not exist
	_ = self.RemoveLocalConfigValue(configdomain.NewParentKey(branch))
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (self *GitTown) RemovePerennialBranchConfiguration() error {
	return self.RemoveLocalConfigValue(configdomain.KeyPerennialBranches)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *GitTown) SetMainBranch(branch domain.LocalBranchName) error {
	self.Config.MainBranch = branch
	return self.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *GitTown) SetNewBranchPush(value configdomain.NewBranchPush, global bool) error {
	setting := strconv.FormatBool(bool(value))
	if global {
		return self.SetGlobalConfigValue(configdomain.KeyPushNewBranches, setting)
	}
	return self.SetLocalConfigValue(configdomain.KeyPushNewBranches, setting)
}

// SetOffline updates whether Git Town is in offline mode.
func (self *GitTown) SetOffline(value configdomain.Offline) error {
	self.Config.Offline = value
	return self.SetGlobalConfigValue(configdomain.KeyOffline, value.String())
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *GitTown) SetParent(branch, parentBranch domain.LocalBranchName) error {
	return self.SetLocalConfigValue(configdomain.NewParentKey(branch), parentBranch.String())
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *GitTown) SetPerennialBranches(branches domain.LocalBranchNames) error {
	self.PerennialBranches = branches
	return self.SetLocalConfigValue(configdomain.KeyPerennialBranches, branches.Join(" "))
}

// SetPushHook updates the configured push-hook strategy.
func (self *GitTown) SetPushHookGlobally(value configdomain.PushHook) error {
	return self.SetGlobalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *GitTown) SetPushHookLocally(value configdomain.PushHook) error {
	return self.SetLocalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
}

// SetShouldShipDeleteRemoteBranch updates the configured delete-remote-branch strategy.
func (self *GitTown) SetShouldShipDeleteRemoteBranch(value configdomain.ShipDeleteTrackingBranch) error {
	return self.SetLocalConfigValue(configdomain.KeyShipDeleteRemoteBranch, strconv.FormatBool(value.Bool()))
}

// SetShouldSyncUpstream updates the configured sync-upstream strategy.
func (self *GitTown) SetShouldSyncUpstream(value configdomain.SyncUpstream) error {
	return self.SetLocalConfigValue(configdomain.KeySyncUpstream, strconv.FormatBool(value.Bool()))
}

func (self *GitTown) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	return self.SetLocalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

func (self *GitTown) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	return self.SetGlobalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *GitTown) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	return self.SetLocalConfigValue(configdomain.KeySyncPerennialStrategy, strategy.String())
}

// SetTestOrigin sets the origin to be used for testing.
func (self *GitTown) SetTestOrigin(value string) error {
	return self.SetLocalConfigValue(configdomain.KeyTestingRemoteURL, value)
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (self *GitTown) ShouldNewBranchPush() (configdomain.NewBranchPush, error) {
	config := self.LocalOrGlobalConfigValue(configdomain.KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	value, err := gohacks.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, configdomain.KeyPushNewBranches, config)
	}
	return configdomain.NewBranchPush(value), nil
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (self *GitTown) ShouldNewBranchPushGlobal() (configdomain.NewBranchPush, error) {
	config := self.GlobalConfigValue(configdomain.KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	boolValue, err := gohacks.ParseBool(config)
	return configdomain.NewBranchPush(boolValue), err
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (self *GitTown) ShouldShipDeleteOriginBranch() (configdomain.ShipDeleteTrackingBranch, error) {
	setting := self.LocalOrGlobalConfigValue(configdomain.KeyShipDeleteRemoteBranch)
	if setting == "" {
		return true, nil
	}
	result, err := strconv.ParseBool(setting)
	if err != nil {
		return true, fmt.Errorf(messages.ValueInvalid, configdomain.KeyShipDeleteRemoteBranch, setting)
	}
	return configdomain.ShipDeleteTrackingBranch(result), nil
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (self *GitTown) ShouldSyncUpstream() (configdomain.SyncUpstream, error) {
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncUpstream)
	if text == "" {
		return true, nil
	}
	boolValue, err := gohacks.ParseBool(text)
	return configdomain.SyncUpstream(boolValue), err
}

// SyncBeforeShip indicates whether a sync should be performed before a ship.
func (self *GitTown) SyncBeforeShip() (configdomain.SyncBeforeShip, error) {
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncBeforeShip)
	if text == "" {
		return false, nil
	}
	boolValue, err := gohacks.ParseBool(text)
	return configdomain.SyncBeforeShip(boolValue), err
}

func (self *GitTown) SyncFeatureStrategy() (configdomain.SyncFeatureStrategy, error) {
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncFeatureStrategy)
	return configdomain.NewSyncFeatureStrategy(text)
}

func (self *GitTown) SyncFeatureStrategyGlobal() (configdomain.SyncFeatureStrategy, error) {
	setting := self.GlobalConfigValue(configdomain.KeySyncFeatureStrategy)
	return configdomain.NewSyncFeatureStrategy(setting)
}

// SyncPerennialStrategy provides the currently configured sync-perennial strategy.
func (self *GitTown) SyncPerennialStrategy() (configdomain.SyncPerennialStrategy, error) {
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncPerennialStrategy)
	return configdomain.NewSyncPerennialStrategy(text)
}
