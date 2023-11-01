package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git/giturl"
	"github.com/git-town/git-town/v10/src/gohacks/slice"
	"github.com/git-town/git-town/v10/src/messages"
)

// GitTown provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type GitTown struct {
	Git
	originURLCache OriginURLCache
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

func (self *GitTown) DeprecatedNewBranchPushFlagGlobal() string {
	return self.config.Global[KeyDeprecatedNewBranchPushFlag]
}

func (self *GitTown) DeprecatedNewBranchPushFlagLocal() string {
	return self.config.Local[KeyDeprecatedNewBranchPushFlag]
}

func (self *GitTown) DeprecatedPushVerifyFlagGlobal() string {
	return self.config.Global[KeyDeprecatedPushVerify]
}

func (self *GitTown) DeprecatedPushVerifyFlagLocal() string {
	return self.config.Local[KeyDeprecatedPushVerify]
}

func DetermineOriginURL(originURL, originOverride string, originURLCache OriginURLCache) *giturl.Parts {
	cached, has := originURLCache[originURL]
	if has {
		return cached
	}
	url := giturl.Parse(originURL)
	if originOverride != "" {
		url.Host = originOverride
	}
	originURLCache[originURL] = url
	return url
}

func NewGitTown(gitConfig GitConfig, runner runner) *GitTown {
	return &GitTown{
		Git:            NewGit(gitConfig, runner),
		originURLCache: OriginURLCache{},
	}
}

// GitAlias provides the currently set alias for the given Git Town command.
func (self *GitTown) GitAlias(alias Alias) string {
	return self.GlobalConfigValue(NewAliasKey(alias))
}

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (self *GitTown) GitHubToken() string {
	return self.LocalOrGlobalConfigValue(KeyGithubToken)
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (self *GitTown) GitLabToken() string {
	return self.LocalOrGlobalConfigValue(KeyGitlabToken)
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (self *GitTown) GiteaToken() string {
	return self.LocalOrGlobalConfigValue(KeyGiteaToken)
}

// HasBranchInformation indicates whether this configuration contains any branch hierarchy entries.
func (self *GitTown) HasBranchInformation() bool {
	for key := range self.config.Local {
		if strings.HasPrefix(key.Name, "git-town-branch.") {
			return true
		}
	}
	return false
}

// HostingService provides the type-safe name of the code hosting connector to use.
// This function caches its result and can be queried repeatedly.
func (self *GitTown) HostingService() (Hosting, error) {
	return NewHosting(self.HostingServiceName())
}

// HostingServiceName provides the name of the code hosting connector to use.
func (self *GitTown) HostingServiceName() string {
	return self.LocalOrGlobalConfigValue(KeyCodeHostingDriver)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *GitTown) IsMainBranch(branch domain.LocalBranchName) bool {
	return branch == self.MainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (self *GitTown) IsOffline() (bool, error) {
	config := self.GlobalConfigValue(KeyOffline)
	if config == "" {
		return false, nil
	}
	result, err := ParseBool(config)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, KeyOffline, config)
	}
	return result, nil
}

// Lineage provides the configured ancestry information for this Git repo.
func (self *GitTown) Lineage(deleteEntry func(Key) error) Lineage {
	lineage := Lineage{}
	for _, key := range self.LocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := domain.NewLocalBranchName(strings.TrimSuffix(strings.TrimPrefix(key.Name, "git-town-branch."), ".parent"))
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
	mainBranch := self.LocalOrGlobalConfigValue(KeyMainBranch)
	if mainBranch == "" {
		return domain.EmptyLocalBranchName()
	}
	return domain.NewLocalBranchName(mainBranch)
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (self *GitTown) OriginOverride() string {
	return self.LocalConfigValue(KeyCodeHostingOriginHostname)
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
	result := self.LocalOrGlobalConfigValue(KeyPerennialBranches)
	if result == "" {
		return domain.LocalBranchNames{}
	}
	return domain.NewLocalBranchNames(strings.Split(result, " ")...)
}

// PullBranchStrategy provides the currently configured pull branch strategy.
func (self *GitTown) PullBranchStrategy() (PullBranchStrategy, error) {
	text := self.LocalOrGlobalConfigValue(KeyPullBranchStrategy)
	return NewPullBranchStrategy(text)
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHook() (bool, error) {
	err := self.updateDeprecatedSetting(KeyDeprecatedPushVerify, KeyPushHook)
	if err != nil {
		return false, err
	}
	setting := self.LocalOrGlobalConfigValue(KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, KeyPushHook, setting)
	}
	return result, nil
}

// PushHook provides the currently configured push-hook setting.
func (self *GitTown) PushHookGlobal() (bool, error) {
	err := self.updateDeprecatedGlobalSetting(KeyDeprecatedPushVerify, KeyPushHook)
	if err != nil {
		return false, err
	}
	setting := self.GlobalConfigValue(KeyPushHook)
	if setting == "" {
		return true, nil
	}
	result, err := ParseBool(setting)
	if err != nil {
		return false, fmt.Errorf(messages.ValueGlobalInvalid, KeyPushHook, setting)
	}
	return result, nil
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
	return nil
}

// RemoveMainBranchConfiguration removes the configuration entry for the main branch name.
func (self *GitTown) RemoveMainBranchConfiguration() error {
	return self.RemoveLocalConfigValue(KeyMainBranch)
}

// RemoveParent removes the parent branch entry for the given branch
// from the Git configuration.
func (self *GitTown) RemoveParent(branch domain.LocalBranchName) {
	// ignoring errors here because the entry might not exist
	_ = self.RemoveLocalConfigValue(NewParentKey(branch))
}

// RemovePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (self *GitTown) RemovePerennialBranchConfiguration() error {
	return self.RemoveLocalConfigValue(KeyPerennialBranches)
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (self *GitTown) SetCodeHostingDriver(value string) error {
	self.config.Local[KeyCodeHostingDriver] = value
	err := self.Run("git", "config", KeyCodeHostingDriver.String(), value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (self *GitTown) SetCodeHostingOriginHostname(value string) error {
	self.config.Local[KeyCodeHostingOriginHostname] = value
	err := self.Run("git", "config", KeyCodeHostingOriginHostname.String(), value)
	return err
}

// SetColorUI configures whether Git output contains color codes.
func (self *GitTown) SetColorUI(value string) error {
	err := self.Run("git", "config", "color.ui", value)
	return err
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *GitTown) SetMainBranch(branch domain.LocalBranchName) error {
	err := self.SetLocalConfigValue(KeyMainBranch, branch.String())
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (self *GitTown) SetNewBranchPush(value bool, global bool) error {
	setting := strconv.FormatBool(value)
	if global {
		err := self.SetGlobalConfigValue(KeyPushNewBranches, setting)
		return err
	}
	err := self.SetLocalConfigValue(KeyPushNewBranches, setting)
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (self *GitTown) SetOffline(value bool) error {
	err := self.SetGlobalConfigValue(KeyOffline, strconv.FormatBool(value))
	return err
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *GitTown) SetParent(branch, parentBranch domain.LocalBranchName) error {
	err := self.SetLocalConfigValue(NewParentKey(branch), parentBranch.String())
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (self *GitTown) SetPerennialBranches(branches domain.LocalBranchNames) error {
	err := self.SetLocalConfigValue(KeyPerennialBranches, branches.Join(" "))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (self *GitTown) SetPullBranchStrategy(strategy PullBranchStrategy) error {
	err := self.SetLocalConfigValue(KeyPullBranchStrategy, strategy.String())
	return err
}

// SetPushHook updates the configured pull branch strategy.
func (self *GitTown) SetPushHookGlobally(value bool) error {
	err := self.SetGlobalConfigValue(KeyPushHook, strconv.FormatBool(value))
	return err
}

// SetPushHookLocally updates the configured pull branch strategy.
func (self *GitTown) SetPushHookLocally(value bool) error {
	err := self.SetLocalConfigValue(KeyPushHook, strconv.FormatBool(value))
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (self *GitTown) SetShouldShipDeleteRemoteBranch(value bool) error {
	err := self.SetLocalConfigValue(KeyShipDeleteRemoteBranch, strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (self *GitTown) SetShouldSyncUpstream(value bool) error {
	err := self.SetLocalConfigValue(KeySyncUpstream, strconv.FormatBool(value))
	return err
}

func (self *GitTown) SetSyncStrategy(value SyncStrategy) error {
	err := self.SetLocalConfigValue(KeySyncStrategy, value.name)
	return err
}

func (self *GitTown) SetSyncStrategyGlobal(value SyncStrategy) error {
	err := self.SetGlobalConfigValue(KeySyncStrategy, value.name)
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (self *GitTown) SetTestOrigin(value string) error {
	err := self.SetLocalConfigValue(KeyTestingRemoteURL, value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (self *GitTown) ShouldNewBranchPush() (bool, error) {
	err := self.updateDeprecatedSetting(KeyDeprecatedNewBranchPushFlag, KeyPushNewBranches)
	if err != nil {
		return false, err
	}
	config := self.LocalOrGlobalConfigValue(KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	value, err := ParseBool(config)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, KeyPushNewBranches, config)
	}
	return value, nil
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (self *GitTown) ShouldNewBranchPushGlobal() (bool, error) {
	err := self.updateDeprecatedGlobalSetting(KeyDeprecatedNewBranchPushFlag, KeyPushNewBranches)
	if err != nil {
		return false, err
	}
	config := self.GlobalConfigValue(KeyPushNewBranches)
	if config == "" {
		return false, nil
	}
	return ParseBool(config)
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (self *GitTown) ShouldShipDeleteOriginBranch() (bool, error) {
	setting := self.LocalOrGlobalConfigValue(KeyShipDeleteRemoteBranch)
	if setting == "" {
		return true, nil
	}
	result, err := strconv.ParseBool(setting)
	if err != nil {
		return true, fmt.Errorf(messages.ValueInvalid, KeyShipDeleteRemoteBranch, setting)
	}
	return result, nil
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (self *GitTown) ShouldSyncUpstream() (bool, error) {
	text := self.LocalOrGlobalConfigValue(KeySyncUpstream)
	if text == "" {
		return true, nil
	}
	return ParseBool(text)
}

func (self *GitTown) SyncStrategy() (SyncStrategy, error) {
	text := self.LocalOrGlobalConfigValue(KeySyncStrategy)
	return ToSyncStrategy(text)
}

func (self *GitTown) SyncStrategyGlobal() (SyncStrategy, error) {
	setting := self.GlobalConfigValue(KeySyncStrategy)
	return ToSyncStrategy(setting)
}

func (self *GitTown) updateDeprecatedGlobalSetting(deprecatedKey, newKey Key) error {
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

func (self *GitTown) updateDeprecatedLocalSetting(deprecatedKey, newKey Key) error {
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

func (self *GitTown) updateDeprecatedSetting(deprecatedKey, newKey Key) error {
	err := self.updateDeprecatedLocalSetting(deprecatedKey, newKey)
	if err != nil {
		return err
	}
	return self.updateDeprecatedGlobalSetting(deprecatedKey, newKey)
}
