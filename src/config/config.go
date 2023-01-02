// Package config provides facilities to read and write the Git Town configuration.
// Git Town stores its configuration in the Git configuration under the prefix "git-town".
// It supports both the Git configuration for the local repository as well as the global Git configuration in `~/.gitconfig`.
// You can manually read the Git configuration entries for Git Town by running `git config --get-regexp git-town`.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/git-town/git-town/v7/src/stringslice"
)

type Config struct {
	gitConfig gitConfig
}

func NewConfiguration(shell run.Shell) Config {
	return Config{
		gitConfig: gitConfig{
			localConfigCache:  loadGitConfig(shell, false),
			globalConfigCache: loadGitConfig(shell, true),
			shell:             shell,
		},
	}
}

// Reload refreshes the cached configuration information.
func (c *Config) Reload() {
	c.gitConfig.Reload()
}

func (c *Config) RemoveLocalGitConfiguration() error {
	return c.gitConfig.RemoveLocalGitConfiguration()
}

func (c *Config) SetLocalConfigValueInTests(key, value string) (*run.Result, error) {
	if flag.Lookup("test.v") == nil {
		panic("this function is only for tests")
	}
	return c.gitConfig.SetLocalConfigValue(key, value)
}

func (c *Config) SetGlobalConfigValueInTests(key, value string) (*run.Result, error) {
	if flag.Lookup("test.v") == nil {
		panic("this function is only for tests")
	}
	return c.gitConfig.SetGlobalConfigValue(key, value)
}

// AddGitAlias sets the given Git alias.
func (c *Config) AddGitAlias(command string) (*run.Result, error) {
	return c.gitConfig.SetGlobalConfigValue("alias."+command, "town "+command)
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (c *Config) AddToPerennialBranches(branchNames ...string) error {
	return c.SetPerennialBranches(append(c.PerennialBranches(), branchNames...))
}

// List provides the names of all parent branches for the given branch.
//
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (c *Config) AncestorBranches(branchName string) []string {
	parentBranchMap := c.ParentBranchMap()
	current := branchName
	result := []string{}
	for {
		if c.IsMainBranch(current) || c.IsPerennialBranch(current) {
			return result
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// IsAncestorBranch indicates whether the given branch is an ancestor of the other given branch.
func (c *Config) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := c.AncestorBranches(branchName)
	return stringslice.Contains(ancestorBranches, ancestorBranchName)
}

// BranchAncestryRoots provides the branches with children and no parents.
func (c *Config) BranchAncestryRoots() []string {
	parentMap := c.ParentBranchMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// ChildBranches provides the names of all branches for which the given branch
// is a parent.
func (c *Config) ChildBranches(branchName string) []string {
	result := []string{}
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := c.gitConfig.localConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return result
}

// HostingService provides the name of the code hosting driver to use.
func (c *Config) HostingService() string {
	return c.localOrGlobalConfigValue("git-town.code-hosting-driver")
}

// OriginOverride provides the override for the origin hostname from the Git Town configuration.
func (c *Config) OriginOverride() string {
	return c.gitConfig.localConfigValue("git-town.code-hosting-origin-hostname")
}

// DeleteMainBranchConfiguration removes the configuration entry for the main branch name.
func (c *Config) DeleteMainBranchConfiguration() error {
	return c.gitConfig.removeLocalConfigValue("git-town.main-branch-name")
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Config) DeleteParentBranch(branchName string) error {
	return c.gitConfig.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// DeletePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (c *Config) DeletePerennialBranchConfiguration() error {
	return c.gitConfig.removeLocalConfigValue("git-town.perennial-branch-names")
}

// GitAlias provides the currently set alias for the given Git Town command.
func (c *Config) GitAlias(command string) string {
	return c.gitConfig.globalConfigValue("alias." + command)
}

// GitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (c *Config) GitHubToken() string {
	return c.localOrGlobalConfigValue("git-town.github-token")
}

// GitLabToken provides the content of the GitLab API token stored in the local or global Git Town configuration.
func (c *Config) GitLabToken() string {
	return c.localOrGlobalConfigValue("git-town.gitlab-token")
}

// GiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (c *Config) GiteaToken() string {
	return c.localOrGlobalConfigValue("git-town.gitea-token")
}

// HasBranchInformation indicates whether this configuration contains any branch hierarchy entries.
func (c *Config) HasBranchInformation() bool {
	for key := range c.gitConfig.localConfigCache {
		if strings.HasPrefix(key, "git-town-branch.") {
			return true
		}
	}
	return false
}

// HasParentBranch returns whether or not the given branch has a parent.
func (c *Config) HasParentBranch(branchName string) bool {
	return c.ParentBranch(branchName) != ""
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (c *Config) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.IsPerennialBranch(branchName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (c *Config) IsMainBranch(branchName string) bool {
	return branchName == c.MainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (c *Config) IsOffline() bool {
	config := c.gitConfig.globalConfigValue("git-town.offline")
	if config == "" {
		return false
	}
	result, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.offline: %q. Please provide either true or false. Considering false for now.", config)
		fmt.Println()
		return false
	}
	return result
}

// IsPerennialBranch indicates whether the branch with the given name is
// a perennial branch.
func (c *Config) IsPerennialBranch(branchName string) bool {
	perennialBranches := c.PerennialBranches()
	return stringslice.Contains(perennialBranches, branchName)
}

// localConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (c *Config) localConfigKeysMatching(toMatch string) []string {
	result := []string{}
	re := regexp.MustCompile(toMatch)
	for key := range c.gitConfig.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// localOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (c *Config) localOrGlobalConfigValue(key string) string {
	local := c.gitConfig.localConfigValue(key)
	if local != "" {
		return local
	}
	return c.gitConfig.globalConfigValue(key)
}

// MainBranch provides the name of the main branch.
func (c *Config) MainBranch() string {
	return c.localOrGlobalConfigValue("git-town.main-branch-name")
}

// ParentBranchMap returns a map from branch name to its parent branch.
func (c *Config) ParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.gitConfig.localConfigValue(key)
		result[child] = parent
	}
	return result
}

// ParentBranch provides the name of the parent branch of the given branch.
func (c *Config) ParentBranch(branchName string) string {
	return c.gitConfig.localConfigValue("git-town-branch." + branchName + ".parent")
}

// PerennialBranches returns all branches that are marked as perennial.
func (c *Config) PerennialBranches() []string {
	result := c.localOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// PullBranchStrategy provides the currently configured pull branch strategy.
func (c *Config) PullBranchStrategy() string {
	config := c.localOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// PushVerify provides the currently configured pull branch strategy.
func (c *Config) PushVerify() bool {
	config := c.localOrGlobalConfigValue("git-town.push-verify")
	if config == "" {
		return true
	}
	result, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.push-verify: %q. Please provide either true or false. Considering true for now.", config)
		fmt.Println()
		return true
	}
	return result
}

// OriginURL provides the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Config) OriginURL() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	res, _ := c.gitConfig.shell.Run("git", "remote", "get-url", "origin")
	return res.OutputSanitized()
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (c *Config) RemoveFromPerennialBranches(branchName string) error {
	return c.SetPerennialBranches(stringslice.Remove(c.PerennialBranches(), branchName))
}

// RemoveGitAlias removes the given Git alias.
func (c *Config) RemoveGitAlias(command string) (*run.Result, error) {
	return c.gitConfig.removeGlobalConfigValue("alias." + command)
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (c *Config) SetCodeHostingDriver(value string) error {
	const key = "git-town.code-hosting-driver"
	c.gitConfig.localConfigCache[key] = value
	_, err := c.gitConfig.shell.Run("git", "config", key, value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (c *Config) SetCodeHostingOriginHostname(value string) error {
	const key = "git-town.code-hosting-origin-hostname"
	c.gitConfig.localConfigCache[key] = value
	_, err := c.gitConfig.shell.Run("git", "config", key, value)
	return err
}

// SetColorUI configures whether Git output contains color codes.
func (c *Config) SetColorUI(value string) error {
	_, err := c.gitConfig.shell.Run("git", "config", "color.ui", value)
	return err
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Config) SetMainBranch(branchName string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.main-branch-name", branchName)
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches to origin.
func (c *Config) SetNewBranchPush(value bool, global bool) error {
	if global {
		_, err := c.gitConfig.SetGlobalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
		return err
	}
	_, err := c.gitConfig.SetLocalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
	return err
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Config) SetParentBranch(branchName, parentBranchName string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (c *Config) SetPerennialBranches(branchNames []string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Config) SetPullBranchStrategy(strategy string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.pull-branch-strategy", strategy)
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Config) SetPushVerify(strategy string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.push-verify", strategy)
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (c *Config) SetOffline(value bool) error {
	_, err := c.gitConfig.SetGlobalConfigValue("git-town.offline", strconv.FormatBool(value))
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (c *Config) SetShouldShipDeleteRemoteBranch(value bool) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.ship-delete-remote-branch", strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (c *Config) SetShouldSyncUpstream(value bool) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.sync-upstream", strconv.FormatBool(value))
	return err
}

func (c *Config) SetSyncStrategy(value string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.sync-strategy", value)
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (c *Config) SetTestOrigin(value string) error {
	_, err := c.gitConfig.SetLocalConfigValue("git-town.testing.remote-url", value)
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to origin.
func (c *Config) ShouldNewBranchPush() bool {
	config := c.localOrGlobalConfigValue("git-town.new-branch-push-flag")
	if config == "" {
		return false
	}
	value, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.new-branch-push-flag: %q. Please provide either true or false. Considering false for now.\n", config)
		return false
	}
	return value
}

// ShouldNewBranchPushGlobal indictes whether the global configuration requires to push
// freshly created branches to origin.
func (c *Config) ShouldNewBranchPushGlobal() bool {
	config := c.gitConfig.globalConfigValue("git-town.new-branch-push-flag")
	return config == "true"
}

// ShouldShipDeleteOriginBranch indicates whether to delete the remote branch after shipping.
func (c *Config) ShouldShipDeleteOriginBranch() bool {
	setting := c.localOrGlobalConfigValue("git-town.ship-delete-remote-branch")
	if setting == "" {
		return true
	}
	result, err := strconv.ParseBool(setting)
	if err != nil {
		fmt.Printf("Invalid value for git-town.ship-delete-remote-branch: %q. Please provide either true or false. Considering true for now.\n", setting)
		return true
	}
	return result
}

// ShouldSyncUpstream indicates whether this repo should sync with its upstream.
func (c *Config) ShouldSyncUpstream() bool {
	return c.localOrGlobalConfigValue("git-town.sync-upstream") != "false"
}

func (c *Config) SyncStrategy() string {
	setting := c.localOrGlobalConfigValue("git-town.sync-strategy")
	if setting == "" {
		setting = "merge"
	}
	return setting
}

// ValidateIsOnline asserts that Git Town is not in offline mode.
func (c *Config) ValidateIsOnline() error {
	if c.IsOffline() {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
