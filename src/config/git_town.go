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
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/src/messages"
)

// GitTown provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type GitTown struct {
	gitconfig.CachedAccess      // access to the Git configuration settings
	DryRun                 bool // whether to dry-run Git commands in this repo
	originURLCache         OriginURLCache
}

type OriginURLCache map[string]*giturl.Parts

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (self *GitTown) AddToPerennialBranches(branches ...domain.LocalBranchName) error {
	return self.SetPerennialBranches(append(self.PerennialBranches(), branches...))
}

func (self *GitTown) BranchTypes() domain.BranchTypes {
	return domain.BranchTypes{
		MainBranch:        self.MainBranch(),
		PerennialBranches: self.PerennialBranches(),
	}
}

func DetermineOriginURL(originURL string, originOverride configdomain.OriginHostnameOverride, originURLCache OriginURLCache) *giturl.Parts {
	cached, has := originURLCache[originURL]
	if has {
		return cached
	}
	url := giturl.Parse(originURL)
	if originOverride != "" {
		url.Host = string(originOverride)
	}
	originURLCache[originURL] = url
	return url
}

func NewGitTown(gitConfig gitconfig.LocalGlobal, runner gitconfig.Runner) *GitTown {
	return &GitTown{
		CachedAccess:   gitconfig.NewGit(gitConfig, runner),
		DryRun:         false,
		originURLCache: OriginURLCache{},
	}
}

// ContainsLineage indicates whether this configuration contains any lineage entries.
func (self *GitTown) ContainsLineage() bool {
	for key := range self.LocalGlobal.Local {
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

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (self *GitTown) GitHubToken() configdomain.GitHubToken {
	return configdomain.GitHubToken(self.LocalOrGlobalConfigValue(configdomain.KeyGithubToken))
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (self *GitTown) GitLabToken() configdomain.GitLabToken {
	return configdomain.GitLabToken(self.LocalOrGlobalConfigValue(configdomain.KeyGitlabToken))
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (self *GitTown) GiteaToken() configdomain.GiteaToken {
	return configdomain.GiteaToken(self.LocalOrGlobalConfigValue(configdomain.KeyGiteaToken))
}

// HostingService provides the type-safe name of the code hosting connector to use.
// This function caches its result and can be queried repeatedly.
func (self *GitTown) HostingService() (configdomain.Hosting, error) {
	return configdomain.NewHosting(self.HostingServiceName())
}

// HostingServiceName provides the name of the code hosting connector to use.
func (self *GitTown) HostingServiceName() string {
	_ = self.updateDeprecatedSetting(configdomain.KeyDeprecatedCodeHostingDriver, configdomain.KeyCodeHostingPlatform)
	return self.LocalOrGlobalConfigValue(configdomain.KeyCodeHostingPlatform)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *GitTown) IsMainBranch(branch domain.LocalBranchName) bool {
	return branch == self.MainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (self *GitTown) IsOffline() (configdomain.Offline, error) {
	config := self.GlobalConfigValue(configdomain.KeyOffline)
	if config == "" {
		return false, nil
	}
	boolValue, err := confighelpers.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, configdomain.KeyOffline, config)
	}
	return configdomain.Offline(boolValue), nil
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

// MainBranch provides the name of the main branch.
func (self *GitTown) MainBranch() domain.LocalBranchName {
	_ = self.updateDeprecatedSetting(configdomain.KeyDeprecatedMainBranchName, configdomain.KeyMainBranch)
	mainBranch := self.LocalOrGlobalConfigValue(configdomain.KeyMainBranch)
	if mainBranch == "" {
		return domain.EmptyLocalBranchName()
	}
	return domain.NewLocalBranchName(mainBranch)
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
	return DetermineOriginURL(text, self.OriginOverride(), self.originURLCache)
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

// PerennialBranches returns all branches that are marked as perennial.
func (self *GitTown) PerennialBranches() domain.LocalBranchNames {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedPerennialBranchNames, configdomain.KeyPerennialBranches)
	if err != nil {
		return domain.NewLocalBranchNames()
	}
	result := self.LocalOrGlobalConfigValue(configdomain.KeyPerennialBranches)
	if result == "" {
		return domain.LocalBranchNames{}
	}
	return domain.NewLocalBranchNames(strings.Split(result, " ")...)
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHook() (configdomain.PushHook, error) {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedPushVerify, configdomain.KeyPushHook)
	if err != nil {
		return false, err
	}
	setting := self.LocalOrGlobalConfigValue(configdomain.KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := confighelpers.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, configdomain.KeyPushHook, setting)
	}
	return configdomain.PushHook(result), nil
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHookGlobal() (configdomain.PushHook, error) {
	err := self.updateDeprecatedGlobalSetting(configdomain.KeyDeprecatedPushVerify, configdomain.KeyPushHook)
	if err != nil {
		return false, err
	}
	setting := self.GlobalConfigValue(configdomain.KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := confighelpers.ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueGlobalInvalid, configdomain.KeyPushHook, setting)
	}
	return configdomain.PushHook(result), nil
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (self *GitTown) RemoveFromPerennialBranches(branch domain.LocalBranchName) error {
	perennialBranches := self.PerennialBranches()
	slice.Remove(&perennialBranches, branch)
	return self.SetPerennialBranches(perennialBranches)
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
	err := self.SetLocalConfigValue(configdomain.KeyMainBranch, branch.String())
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *GitTown) SetNewBranchPush(value configdomain.NewBranchPush, global bool) error {
	setting := strconv.FormatBool(bool(value))
	if global {
		err := self.SetGlobalConfigValue(configdomain.KeyPushNewBranches, setting)
		return err
	}
	err := self.SetLocalConfigValue(configdomain.KeyPushNewBranches, setting)
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (self *GitTown) SetOffline(value configdomain.Offline) error {
	err := self.SetGlobalConfigValue(configdomain.KeyOffline, strconv.FormatBool(value.Bool()))
	return err
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *GitTown) SetParent(branch, parentBranch domain.LocalBranchName) error {
	err := self.SetLocalConfigValue(configdomain.NewParentKey(branch), parentBranch.String())
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *GitTown) SetPerennialBranches(branches domain.LocalBranchNames) error {
	err := self.SetLocalConfigValue(configdomain.KeyPerennialBranches, branches.Join(" "))
	return err
}

// SetPushHook updates the configured push-hook strategy.
func (self *GitTown) SetPushHookGlobally(value configdomain.PushHook) error {
	err := self.SetGlobalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
	return err
}

// SetPushHookLocally updates the locally configured push-hook strategy.
func (self *GitTown) SetPushHookLocally(value configdomain.PushHook) error {
	err := self.SetLocalConfigValue(configdomain.KeyPushHook, strconv.FormatBool(bool(value)))
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured delete-remote-branch strategy.
func (self *GitTown) SetShouldShipDeleteRemoteBranch(value configdomain.ShipDeleteTrackingBranch) error {
	err := self.SetLocalConfigValue(configdomain.KeyShipDeleteRemoteBranch, strconv.FormatBool(value.Bool()))
	return err
}

// SetShouldSyncUpstream updates the configured sync-upstream strategy.
func (self *GitTown) SetShouldSyncUpstream(value configdomain.SyncUpstream) error {
	err := self.SetLocalConfigValue(configdomain.KeySyncUpstream, strconv.FormatBool(value.Bool()))
	return err
}

func (self *GitTown) SetSyncFeatureStrategy(value configdomain.SyncFeatureStrategy) error {
	err := self.SetLocalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
	return err
}

func (self *GitTown) SetSyncFeatureStrategyGlobal(value configdomain.SyncFeatureStrategy) error {
	err := self.SetGlobalConfigValue(configdomain.KeySyncFeatureStrategy, value.Name)
	return err
}

// SetSyncPerennialStrategy updates the configured sync-perennial strategy.
func (self *GitTown) SetSyncPerennialStrategy(strategy configdomain.SyncPerennialStrategy) error {
	err := self.SetLocalConfigValue(configdomain.KeySyncPerennialStrategy, strategy.String())
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (self *GitTown) SetTestOrigin(value string) error {
	err := self.SetLocalConfigValue(configdomain.KeyTestingRemoteURL, value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (self *GitTown) ShouldNewBranchPush() (configdomain.NewBranchPush, error) {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedNewBranchPushFlag, configdomain.KeyPushNewBranches)
	if err != nil {
		return false, err
	}
	config := self.LocalOrGlobalConfigValue(configdomain.KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	value, err := confighelpers.ParseBool(config)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, configdomain.KeyPushNewBranches, config)
	}
	return configdomain.NewBranchPush(value), nil
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (self *GitTown) ShouldNewBranchPushGlobal() (configdomain.NewBranchPush, error) {
	err := self.updateDeprecatedGlobalSetting(configdomain.KeyDeprecatedNewBranchPushFlag, configdomain.KeyPushNewBranches)
	if err != nil {
		return false, err
	}
	config := self.GlobalConfigValue(configdomain.KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	boolValue, err := confighelpers.ParseBool(config)
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
	boolValue, err := confighelpers.ParseBool(text)
	return configdomain.SyncUpstream(boolValue), err
}

// SyncBeforeShip indicates whether a sync should be performed before a ship.
func (self *GitTown) SyncBeforeShip() (configdomain.SyncBeforeShip, error) {
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncBeforeShip)
	if text == "" {
		return false, nil
	}
	boolValue, err := confighelpers.ParseBool(text)
	return configdomain.SyncBeforeShip(boolValue), err
}

func (self *GitTown) SyncFeatureStrategy() (configdomain.SyncFeatureStrategy, error) {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedSyncStrategy, configdomain.KeySyncFeatureStrategy)
	if err != nil {
		return configdomain.SyncFeatureStrategyMerge, err
	}
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncFeatureStrategy)
	return configdomain.NewSyncFeatureStrategy(text)
}

func (self *GitTown) SyncFeatureStrategyGlobal() (configdomain.SyncFeatureStrategy, error) {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedSyncStrategy, configdomain.KeySyncFeatureStrategy)
	if err != nil {
		return configdomain.SyncFeatureStrategyMerge, err
	}
	setting := self.GlobalConfigValue(configdomain.KeySyncFeatureStrategy)
	return configdomain.NewSyncFeatureStrategy(setting)
}

// SyncPerennialStrategy provides the currently configured sync-perennial strategy.
func (self *GitTown) SyncPerennialStrategy() (configdomain.SyncPerennialStrategy, error) {
	err := self.updateDeprecatedSetting(configdomain.KeyDeprecatedPullBranchStrategy, configdomain.KeySyncPerennialStrategy)
	if err != nil {
		return configdomain.SyncPerennialStrategyRebase, err
	}
	text := self.LocalOrGlobalConfigValue(configdomain.KeySyncPerennialStrategy)
	return configdomain.NewSyncPerennialStrategy(text)
}

func (self *GitTown) updateDeprecatedGlobalSetting(deprecatedKey, newKey configdomain.Key) error {
	deprecatedSetting := self.GlobalConfigValue(deprecatedKey)
	if deprecatedSetting != "" {
		fmt.Printf("I found the deprecated global setting %q.\n", deprecatedKey)
		fmt.Printf("I am upgrading this setting to the new format %q.\n", newKey)
		err := self.RemoveGlobalConfigValue(deprecatedKey)
		if err != nil {
			return err
		}
		err = self.SetGlobalConfigValue(newKey, deprecatedSetting)
		return err
	}
	return nil
}

func (self *GitTown) updateDeprecatedLocalSetting(deprecatedKey, newKey configdomain.Key) error {
	deprecatedSetting := self.LocalConfigValue(deprecatedKey)
	if deprecatedSetting != "" {
		fmt.Printf("I found the deprecated local setting %q.\n", deprecatedKey)
		fmt.Printf("I am upgrading this setting to the new format %q.\n", newKey)
		err := self.RemoveLocalConfigValue(deprecatedKey)
		if err != nil {
			return err
		}
		err = self.SetLocalConfigValue(newKey, deprecatedSetting)
		return err
	}
	return nil
}

func (self *GitTown) updateDeprecatedSetting(deprecatedKey, newKey configdomain.Key) error {
	err := self.updateDeprecatedLocalSetting(deprecatedKey, newKey)
	if err != nil {
		return err
	}
	return self.updateDeprecatedGlobalSetting(deprecatedKey, newKey)
}
