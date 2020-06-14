package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/stringslice"
)

// Config manages the Git Town configuration
// stored in Git metadata in the given local repo and the global Git configuration.
// This class manages which config values are stored in local vs global settings.
type Config struct {

	// localConfigCache is a cache of the Git configuration in the local Git repo.
	localConfigCache map[string]string

	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache map[string]string

	// for running shell commands
	shell command.Shell
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewConfiguration(shell command.Shell) *Config {
	return &Config{
		shell:             shell,
		localConfigCache:  loadGitConfig(shell, false),
		globalConfigCache: loadGitConfig(shell, true),
	}
}

// loadGitConfig provides the Git configuration from the given directory or the global one if the global flag is set.
func loadGitConfig(shell command.Shell, global bool) map[string]string {
	result := map[string]string{}
	cmdArgs := []string{"config", "-lz"}
	if global {
		cmdArgs = append(cmdArgs, "--global")
	} else {
		cmdArgs = append(cmdArgs, "--local")
	}
	res, err := shell.Run("git", cmdArgs...)
	if err != nil {
		return result
	}
	output := res.Output()
	if output == "" {
		return result
	}
	for _, line := range strings.Split(output, "\x00") {
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "\n", 2)
		key, value := parts[0], parts[1]
		result[key] = value
	}
	return result
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (c *Config) AddToPerennialBranches(branchNames ...string) error {
	return c.SetPerennialBranches(append(c.GetPerennialBranches(), branchNames...))
}

// AddGitAlias sets the given Git alias.
func (c *Config) AddGitAlias(command string) (*command.Result, error) {
	return c.SetGlobalConfigValue("alias."+command, "town "+command)
}

// DeleteMainBranchConfiguration removes the configuration entry for the main branch name.
func (c *Config) DeleteMainBranchConfiguration() error {
	return c.removeLocalConfigValue("git-town.main-branch-name")
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Config) DeleteParentBranch(branchName string) error {
	return c.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// DeletePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (c *Config) DeletePerennialBranchConfiguration() error {
	return c.removeLocalConfigValue("git-town.perennial-branch-names")
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (c *Config) GetAncestorBranches(branchName string) (result []string) {
	parentBranchMap := c.GetParentBranchMap()
	current := branchName
	for {
		if c.IsMainBranch(current) || c.IsPerennialBranch(current) {
			return
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// GetBranchAncestryRoots provides the branches with children and no parents.
func (c *Config) GetBranchAncestryRoots() []string {
	parentMap := c.GetParentBranchMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func (c *Config) GetChildBranches(branchName string) (result []string) {
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := c.getLocalConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetCodeHostingDriverName provides the name of the code hosting driver to use.
func (c *Config) GetCodeHostingDriverName() string {
	return c.getLocalOrGlobalConfigValue("git-town.code-hosting-driver")
}

// GetCodeHostingOriginHostname provides the host name of the code hosting server.
func (c *Config) GetCodeHostingOriginHostname() string {
	return c.getLocalConfigValue("git-town.code-hosting-origin-hostname")
}

// getGlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Config) getGlobalConfigValue(key string) string {
	return c.globalConfigCache[key]
}

// getLocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Config) getLocalConfigValue(key string) string {
	return c.localConfigCache[key]
}

// getLocalOrGlobalConfigValue provides the configuration value with the given key from the local and global Git configuration.
// Local configuration takes precedence.
func (c *Config) getLocalOrGlobalConfigValue(key string) string {
	local := c.getLocalConfigValue(key)
	if local != "" {
		return local
	}
	return c.getGlobalConfigValue(key)
}

// GetParentBranchMap returns a map from branch name to its parent branch.
func (c *Config) GetParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range c.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.getLocalConfigValue(key)
		result[child] = parent
	}
	return result
}

// GetGitAlias provides the currently set alias for the given Git Town command.
func (c *Config) GetGitAlias(command string) string {
	return c.getGlobalConfigValue("alias." + command)
}

// GetGitHubToken provides the content of the GitHub API token stored in the local or global Git Town configuration.
func (c *Config) GetGitHubToken() string {
	return c.getLocalOrGlobalConfigValue("git-town.github-token")
}

// GetGiteaToken provides the content of the Gitea API token stored in the local or global Git Town configuration.
func (c *Config) GetGiteaToken() string {
	return c.getLocalOrGlobalConfigValue("git-town.gitea-token")
}

// GetMainBranch returns the name of the main branch.
func (c *Config) GetMainBranch() string {
	return c.getLocalOrGlobalConfigValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func (c *Config) GetParentBranch(branchName string) string {
	return c.getLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func (c *Config) GetPerennialBranches() []string {
	result := c.getLocalOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func (c *Config) GetPullBranchStrategy() string {
	config := c.getLocalOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Config) GetRemoteOriginURL() string {
	remote := os.Getenv("GIT_TOWN_REMOTE")
	if remote != "" {
		return remote
	}
	res, _ := c.shell.Run("git", "remote", "get-url", "origin")
	return res.OutputSanitized()
}

// HasBranchInformation indicates whether this configuration contains any branch hierarchy entries.
func (c *Config) HasBranchInformation() bool {
	for key := range c.localConfigCache {
		if strings.HasPrefix(key, "git-town-branch.") {
			return true
		}
	}
	return false
}

// HasParentBranch returns whether or not the given branch has a parent.
func (c *Config) HasParentBranch(branchName string) bool {
	return c.GetParentBranch(branchName) != ""
}

// IsAncestorBranch indicates whether the given branch is an ancestor of the other given branch.
func (c *Config) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := c.GetAncestorBranches(branchName)
	return stringslice.Contains(ancestorBranches, ancestorBranchName)
}

// IsFeatureBranch indicates whether the branch with the given name is
// a feature branch.
func (c *Config) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.IsPerennialBranch(branchName)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (c *Config) IsMainBranch(branchName string) bool {
	return branchName == c.GetMainBranch()
}

// IsOffline indicates whether Git Town is currently in offline mode.
func (c *Config) IsOffline() bool {
	config := c.getGlobalConfigValue("git-town.offline")
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
	perennialBranches := c.GetPerennialBranches()
	return stringslice.Contains(perennialBranches, branchName)
}

// localConfigKeysMatching provides the names of the Git Town configuration keys matching the given RegExp string.
func (c *Config) localConfigKeysMatching(toMatch string) (result []string) {
	re := regexp.MustCompile(toMatch)
	for key := range c.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

// Reload refreshes the cached configuration information.
func (c *Config) Reload() {
	c.localConfigCache = loadGitConfig(c.shell, false)
	c.globalConfigCache = loadGitConfig(c.shell, true)
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (c *Config) RemoveFromPerennialBranches(branchName string) error {
	return c.SetPerennialBranches(stringslice.Remove(c.GetPerennialBranches(), branchName))
}

// RemoveGitAlias removes the given Git alias.
func (c *Config) RemoveGitAlias(command string) (*command.Result, error) {
	return c.removeGlobalConfigValue("alias." + command)
}

func (c *Config) removeGlobalConfigValue(key string) (*command.Result, error) {
	delete(c.globalConfigCache, key)
	return c.shell.Run("git", "config", "--global", "--unset", key)
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *Config) removeLocalConfigValue(key string) error {
	delete(c.localConfigCache, key)
	_, err := c.shell.Run("git", "config", "--unset", key)
	return err
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (c *Config) RemoveLocalGitConfiguration() error {
	_, err := c.shell.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 128 {
			// Git returns exit code 128 when trying to delete a non-existing config section.
			// This is not an error condition in this workflow so we can ignore it here.
			return nil
		}
		return fmt.Errorf("unexpected error while removing the 'git-town' section from the Git configuration: %v", err)
	}
	return nil
}

// SetCodeHostingDriver sets the "github.code-hosting-driver" setting.
func (c *Config) SetCodeHostingDriver(value string) error {
	const key = "git-town.code-hosting-driver"
	c.localConfigCache[key] = value
	_, err := c.shell.Run("git", "config", key, value)
	return err
}

// SetCodeHostingOriginHostname sets the "github.code-hosting-driver" setting.
func (c *Config) SetCodeHostingOriginHostname(value string) error {
	const key = "git-town.code-hosting-origin-hostname"
	c.localConfigCache[key] = value
	_, err := c.shell.Run("git", "config", key, value)
	return err
}

// SetColorUI configures whether Git output contains color codes.
func (c *Config) SetColorUI(value string) error {
	_, err := c.shell.Run("git", "config", "color.ui", value)
	return err
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (c *Config) SetGlobalConfigValue(key, value string) (*command.Result, error) {
	c.globalConfigCache[key] = value
	return c.shell.Run("git", "config", "--global", key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (c *Config) SetLocalConfigValue(key, value string) (*command.Result, error) {
	c.localConfigCache[key] = value
	return c.shell.Run("git", "config", key, value)
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Config) SetMainBranch(branchName string) error {
	_, err := c.SetLocalConfigValue("git-town.main-branch-name", branchName)
	return err
}

// SetNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Config) SetNewBranchPush(value bool, global bool) error {
	if global {
		_, err := c.SetGlobalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
		return err
	}
	_, err := c.SetLocalConfigValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
	return err
}

// SetOffline updates whether Git Town is in offline mode.
func (c *Config) SetOffline(value bool) error {
	_, err := c.SetGlobalConfigValue("git-town.offline", strconv.FormatBool(value))
	return err
}

// SetTestOrigin sets the origin to be used for testing.
func (c *Config) SetTestOrigin(value string) error {
	_, err := c.SetLocalConfigValue("git-town.testing.remote-url", value)
	return err
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Config) SetParentBranch(branchName, parentBranchName string) error {
	_, err := c.SetLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
	return err
}

// SetPerennialBranches marks the given branches as perennial branches.
func (c *Config) SetPerennialBranches(branchNames []string) error {
	_, err := c.SetLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
	return err
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Config) SetPullBranchStrategy(strategy string) error {
	_, err := c.SetLocalConfigValue("git-town.pull-branch-strategy", strategy)
	return err
}

// SetShouldShipDeleteRemoteBranch updates the configured pull branch strategy.
func (c *Config) SetShouldShipDeleteRemoteBranch(value bool) error {
	_, err := c.SetLocalConfigValue("git-town.ship-delete-remote-branch", strconv.FormatBool(value))
	return err
}

// SetShouldSyncUpstream updates the configured pull branch strategy.
func (c *Config) SetShouldSyncUpstream(value bool) error {
	_, err := c.SetLocalConfigValue("git-town.sync-upstream", strconv.FormatBool(value))
	return err
}

// ShouldNewBranchPush indicates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Config) ShouldNewBranchPush() bool {
	config := c.getLocalOrGlobalConfigValue("git-town.new-branch-push-flag")
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
// freshly created branches up to the origin remote.
func (c *Config) ShouldNewBranchPushGlobal() bool {
	config := c.getGlobalConfigValue("git-town.new-branch-push-flag")
	return config == "true"
}

// ShouldShipDeleteRemoteBranch indicates whether to delete the remote branch after shipping.
func (c *Config) ShouldShipDeleteRemoteBranch() bool {
	setting := c.getLocalOrGlobalConfigValue("git-town.ship-delete-remote-branch")
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
	return c.getLocalOrGlobalConfigValue("git-town.sync-upstream") != "false"
}

// ValidateIsOnline asserts that Git Town is not in offline mode.
func (c *Config) ValidateIsOnline() error {
	if c.IsOffline() {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
