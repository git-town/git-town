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

// Config provides access to the Git Town configuration.
var Config *Configuration

// Configuration manages the Git Town configuration,
// stored in Git metadata in the local repo and in the global Git configuration.
type Configuration struct {

	// localConfig is a cache of the Git configuration in the local directory.
	localConfig ConfigMap

	// globalConfig is a cache of the global Git configuration
	globalConfig ConfigMap
}

// NewConfiguration provides Configuration instances.
func NewConfiguration() *Configuration {
	return &Configuration{localConfig: NewConfigMap(false), globalConfig: NewConfigMap(true)}
}

// AddToPerennialBranches adds the given branch as a perennial branch
func (c *Configuration) AddToPerennialBranches(branchName string) {
	c.SetPerennialBranches(append(c.GetPerennialBranches(), branchName))
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func (c *Configuration) DeleteParentBranch(branchName string) {
	c.removeLocalConfigurationValue("git-town-branch." + branchName + ".parent")
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
	for _, key := range c.getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := c.GetConfigurationValue(key)
		result[child] = parent
	}
	return result
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func (c *Configuration) GetChildBranches(branchName string) (result []string) {
	for _, key := range c.getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := c.GetConfigurationValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return
}

// GetConfigurationValue returns the given configuration value,
// from either global or local Git configuration
// TODO: replace this with the proper named method?
func (c *Configuration) GetConfigurationValue(key string) (result string) {
	return c.localConfig.Get(key)
}

// GetGlobalConfigurationValue returns the global git configuration value for the given key
func (c *Configuration) GetGlobalConfigurationValue(key string) string {
	return c.globalConfig.Get(key)
}

// GetMainBranch returns the name of the main branch.
func (c *Configuration) GetMainBranch() string {
	return c.GetConfigurationValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func (c *Configuration) GetParentBranch(branchName string) string {
	return c.GetConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func (c *Configuration) GetPerennialBranches() []string {
	result := c.GetConfigurationValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func (c *Configuration) GetPullBranchStrategy() string {
	return c.getConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func (c *Configuration) GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := c.GetConfigurationValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.New("git", "remote", "get-url", "origin").Output()
}

// GetRemoteUpstreamURL returns the URL of the "upstream" remote.
func (c *Configuration) GetRemoteUpstreamURL() string {
	return command.New("git", "remote", "get-url", "upstream").Output()
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
func (c *Configuration) HasGlobalConfigurationValue(key string) bool {
	return command.New("git", "config", "-l", "--global", "--name").OutputContainsLine(key)
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

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func (c *Configuration) HasRemote(name string) bool {
	return util.DoesStringArrayContain(c.getRemotes(), name)
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
func (c *Configuration) RemoveAllConfiguration() {
	command.New("git", "config", "--remove-section", "git-town").Output()
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
	c.setConfigurationValue("git-town.main-branch-name", branchName)
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (c *Configuration) SetParentBranch(branchName, parentBranchName string) {
	c.setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func (c *Configuration) SetPerennialBranches(branchNames []string) {
	c.setConfigurationValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func (c *Configuration) SetPullBranchStrategy(strategy string) {
	c.setConfigurationValue("git-town.pull-branch-strategy", strategy)
}

// ShouldNewBranchPush returns whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) ShouldNewBranchPush() bool {
	return util.StringToBool(c.getConfigurationValueWithDefault("git-town.new-branch-push-flag", "false"))
}

// GetGlobalNewBranchPushFlag returns the global configuration for to push
// freshly created branches up to the origin remote.
func (c *Configuration) GetGlobalNewBranchPushFlag() string {
	return c.getGlobalConfigurationValueWithDefault("git-town.new-branch-push-flag", "false")
}

// UpdateOffline updates whether Git Town is in offline mode
func (c *Configuration) UpdateOffline(value bool) {
	c.setGlobalConfigurationValue("git-town.offline", strconv.FormatBool(value))
}

// UpdateShouldNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func (c *Configuration) UpdateShouldNewBranchPush(value bool) {
	c.setConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// UpdateGlobalShouldNewBranchPush updates global whether to push
// freshly created branches up to the origin remote.
func (c *Configuration) UpdateGlobalShouldNewBranchPush(value bool) {
	c.setGlobalConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// Helpers

func (c *Configuration) getGlobalConfigurationValueWithDefault(key, defaultValue string) string {
	value := c.GetGlobalConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Configuration) getConfigurationValueWithDefault(key, defaultValue string) string {
	value := c.GetConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Configuration) getConfigurationKeysMatching(toMatch string) (result []string) {
	re, err := regexp.Compile(toMatch)
	exit.IfWrapf(err, "Error compiling configuration regular expression (%s): %v", toMatch, err)
	return c.localConfig.KeysMatching(re)
}

func (c *Configuration) setConfigurationValue(key, value string) {
	command.New("git", "config", key, value).Run()
	c.localConfig.Set(key, value)
}

func (c *Configuration) setGlobalConfigurationValue(key, value string) {
	command.New("git", "config", "--global", key, value).Run()
	c.globalConfig.Set(key, value)
	c.localConfig.Reset() // Need to reset config in case it was inheriting
}

// removeLocalConfigurationValue deletes the configuration value with the given key from the local Git Town configuration.
func (c *Configuration) removeLocalConfigurationValue(key string) {
	command.New("git", "config", "--unset", key).Run()
	c.localConfig.Delete(key)
}

// Remotes are cached in order to minimize the number of git commands run
var remotes []string
var remotesInitialized bool

func (c *Configuration) getRemotes() []string {
	if !remotesInitialized {
		remotes = command.New("git", "remote").OutputLines()
		remotesInitialized = true
	}
	return remotes
}

// Init

func init() {
	Config = NewConfiguration()
}
