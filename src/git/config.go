/*
This file contains functionality around storing configuration settings
inside Git's metadata storage for the repository.
*/

package git

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

var configMap *ConfigMap
var globalConfigMap *ConfigMap

// Configuration manages the Git Town configuration
// stored in Git metadata in the given local repo and the global Git configuration.
// This class is aware which config values are stored in local vs global settings.
type Configuration struct {

	// localDir contains the directory of the local Git repo.
	localDir string

	// localConfigCache is a cache of the Git configuration in the local directory.
	localConfigCache map[string]string

	// globalConfigCache is a cache of the global Git configuration.
	globalConfigCache map[string]string
}

// Config provides the current configuration.
func Config() *Configuration {
	if currentDirConfig == nil {
		currentDirConfig = NewConfiguration("")
	}
	return currentDirConfig
}

// currentDirConfig provides access to the Git Town configuration in the current working directory.
var currentDirConfig *Configuration

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewConfiguration(dir string) *Configuration {
	result := &Configuration{
		localDir:          dir,
		localConfigCache:  loadCache(dir, false),
		globalConfigCache: loadCache(dir, true),
	}
	return result
}

func loadCache(dir string, global bool) map[string]string {
	result := map[string]string{}
	cmdArgs := []string{"config", "-lz"}
	var res *command.Result
	if global {
		cmdArgs = append(cmdArgs, "--global")
		res = command.RunInDir(dir, "git", cmdArgs...)
	} else {
		cmdArgs = append(cmdArgs, "--local")
		res = command.RunInDir(dir, "git", cmdArgs...)
	}
	if res.Err() != nil && strings.Contains(res.OutputSanitized(), "No such file or directory") {
		return result
	}
	exit.If(res.Err())
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

// getLocalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getLocalConfigValue(key string) string {
	return c.localConfigCache[key]
}

// getGlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getGlobalConfigValue(key string) string {
	return c.globalConfigCache[key]
}

// getGlobalConfigValue provides the configuration value with the given key from the local Git configuration.
func (c *Configuration) getLocalOrGlobalConfigValue(key string) string {
	local := c.getLocalConfigValue(key)
	if local != "" {
		return local
	}
	return c.getGlobalConfigValue(key)
}

// setConfigurationValue sets the local configuration with the given key to the given value.
func (c *Configuration) setLocalConfigValue(key, value string) {
	command.RunInDir(c.localDir, "git", "config", key, value)
	c.localConfigCache[key] = value
}

func (c *Configuration) setGlobalConfigValue(key, value string) {
	command.RunInDir(c.localDir, "git", "config", "--global", key, value)
	c.globalConfigCache[key] = value
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *Configuration) removeLocalConfigValue(key string) {
	command.RunInDir(c.localDir, "git", "config", "--unset", key)
	delete(c.localConfigCache, key)
}

// AddToPerennialBranches adds the given branch as a perennial branch
func (c *Configuration) AddToPerennialBranches(branchName string) {
	c.SetPerennialBranches(append(c.GetPerennialBranches(), branchName))
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Configuration) DeleteParentBranch(branchName string) {
	c.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// EnsureIsFeatureBranch asserts that the given branch is a feature branch.
func (c *Configuration) EnsureIsFeatureBranch(branchName, errorSuffix string) {
	util.Ensure(c.IsFeatureBranch(branchName), fmt.Sprintf("The branch '%s' is not a feature branch. %s", branchName, errorSuffix))
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (c *Configuration) GetAncestorBranches(branchName string) (result []string) {
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

// GetParentBranchMap returns a map from branch name to its parent branch
func (c *Configuration) GetParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range c.getLocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.getLocalConfigValue(key)
		result[child] = parent
	}
	return result
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func (c *Configuration) GetChildBranches(branchName string) (result []string) {
	for _, key := range c.getLocalConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := c.getLocalConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetConfigurationValue returns the given configuration value,
// from either global or local Git configuration
func GetConfigurationValue(key string) (result string) {
	return configMap.Get(key)
}

// GetGlobalConfigurationValue returns the global git configuration value for the given key
func GetGlobalConfigurationValue(key string) string {
	return globalConfigMap.Get(key)
}

// GetMainBranch returns the name of the main branch.
func (c *Configuration) GetMainBranch() string {
	return c.getLocalOrGlobalConfigValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func (c *Configuration) GetParentBranch(branchName string) string {
	return c.getLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func (c *Configuration) GetPerennialBranches() []string {
	result := c.getLocalOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func (c *Configuration) GetPullBranchStrategy() string {
	config := c.getLocalOrGlobalConfigValue("git-town.pull-branch-strategy")
	if config != "" {
		return config
	}
	return "rebase"
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Configuration) GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := c.getLocalConfigValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.RunInDir(c.localDir, "git", "remote", "get-url", "origin").OutputSanitized()
}

// GetURLHostname returns the hostname contained within the given Git URL.
func (c *Configuration) GetURLHostname(url string) string {
	hostnameRegex := regexp.MustCompile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func (c *Configuration) GetURLRepositoryName(url string) string {
	hostname := c.GetURLHostname(url)
	repositoryNameRegex, err := regexp.Compile(".*" + hostname + "[/:](.+)")
	exit.IfWrap(err, "Error compiling repository name regular expression")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}

// HasGlobalConfigurationValue returns whether there is a global git configuration for the given key
func HasGlobalConfigurationValue(key string) bool {
	return command.Run("git", "config", "-l", "--global", "--name").OutputContainsLine(key)
}

// HasParentBranch returns whether or not the given branch has a parent
func (c *Configuration) HasParentBranch(branchName string) bool {
	return c.GetParentBranch(branchName) != ""
}

// IsAncestorBranch returns whether the given branch is an ancestor of the other given branch.
func (c *Configuration) IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := c.GetAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

// IsFeatureBranch returns whether the branch with the given name is
// a feature branch.
func (c *Configuration) IsFeatureBranch(branchName string) bool {
	return !c.IsMainBranch(branchName) && !c.IsPerennialBranch(branchName)
}

// IsMainBranch returns whether the branch with the given name
// is the main branch of the repository.
func (c *Configuration) IsMainBranch(branchName string) bool {
	return branchName == c.GetMainBranch()
}

// IsPerennialBranch returns whether the branch with the given name is
// a perennial branch.
func (c *Configuration) IsPerennialBranch(branchName string) bool {
	perennialBranches := c.GetPerennialBranches()
	return util.DoesStringArrayContain(perennialBranches, branchName)
}

// RemoveAllConfiguration removes all Git Town configuration
func RemoveAllConfiguration() {
	command.Run("git", "config", "--remove-section", "git-town").OutputSanitized()
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration
func (c *Configuration) RemoveOutdatedConfiguration() {
	for child, parent := range c.GetParentBranchMap() {
		if !HasBranch(child) || !HasBranch(parent) {
			c.DeleteParentBranch(child)
		}
	}
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch
func (c *Configuration) RemoveFromPerennialBranches(branchName string) {
	c.SetPerennialBranches(util.RemoveStringFromSlice(c.GetPerennialBranches(), branchName))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (c *Configuration) SetMainBranch(branchName string) {
	c.setLocalConfigValue("git-town.main-branch-name", branchName)
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Configuration) SetParentBranch(branchName, parentBranchName string) {
	c.setLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func (c *Configuration) SetPerennialBranches(branchNames []string) {
	c.setLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Configuration) SetPullBranchStrategy(strategy string) {
	c.setLocalConfigValue("git-town.pull-branch-strategy", strategy)
}

// ShouldNewBranchPush returns whether the current repository is configured to push
// freshly created branches up to the origin remote.
func ShouldNewBranchPush() bool {
	return util.StringToBool(getConfigurationValueWithDefault("git-town.new-branch-push-flag", "false"))
}

// GetGlobalNewBranchPushFlag returns the global configuration for to push
// freshly created branches up to the origin remote.
func GetGlobalNewBranchPushFlag() string {
	return getGlobalConfigurationValueWithDefault("git-town.new-branch-push-flag", "false")
}

// UpdateOffline updates whether Git Town is in offline mode
func UpdateOffline(value bool) {
	setGlobalConfigurationValue("git-town.offline", strconv.FormatBool(value))
}

// UpdateShouldNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func UpdateShouldNewBranchPush(value bool) {
	setConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// UpdateGlobalShouldNewBranchPush updates global whether to push
// freshly created branches up to the origin remote.
func UpdateGlobalShouldNewBranchPush(value bool) {
	setGlobalConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// Helpers

func getGlobalConfigurationValueWithDefault(key, defaultValue string) string {
	value := GetGlobalConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getConfigurationValueWithDefault(key, defaultValue string) string {
	value := GetConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Configuration) getLocalConfigKeysMatching(toMatch string) (result []string) {
	re := regexp.MustCompile(toMatch)
	for key := range c.localConfigCache {
		if re.MatchString(key) {
			result = append(result, key)
		}
	}
	return result
}

func setConfigurationValue(key, value string) {
	command.Run("git", "config", key, value)
	configMap.Set(key, value)
}

func setGlobalConfigurationValue(key, value string) {
	command.Run("git", "config", "--global", key, value)
	globalConfigMap.Set(key, value)
	configMap.Reset() // Need to reset config in case it was inheriting
}

func removeConfigurationValue(key string) {
	command.Run("git", "config", "--unset", key)
	configMap.Delete(key)
}

// Init

func init() {
	configMap = NewConfigMap(false)
	globalConfigMap = NewConfigMap(true)
}
