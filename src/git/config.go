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

// AddToPerennialBranches adds the given branch as a perennial branch
func AddToPerennialBranches(branchName string) {
	SetPerennialBranches(append(GetPerennialBranches(), branchName))
}

// DeleteParentBranch removes the parent branch entry for the given branch
// from the Git configuration.
func DeleteParentBranch(branchName string) {
	removeConfigurationValue("git-town-branch." + branchName + ".parent")
}

// EnsureIsFeatureBranch asserts that the given branch is a feature branch.
func EnsureIsFeatureBranch(branchName, errorSuffix string) {
	util.Ensure(IsFeatureBranch(branchName), fmt.Sprintf("The branch '%s' is not a feature branch. %s", branchName, errorSuffix))
}

// GetAncestorBranches returns the names of all parent branches for the given branch,
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func GetAncestorBranches(branchName string) (result []string) {
	parentBranchMap := GetParentBranchMap()
	current := branchName
	for {
		if IsMainBranch(current) || IsPerennialBranch(current) {
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
func GetParentBranchMap() map[string]string {
	result := map[string]string{}
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := GetConfigurationValue(key)
		result[child] = parent
	}
	return result
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func GetChildBranches(branchName string) (result []string) {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := GetConfigurationValue(key)
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
func GetMainBranch() string {
	return GetConfigurationValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func GetParentBranch(branchName string) string {
	return GetConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func GetPerennialBranches() []string {
	result := GetConfigurationValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func GetPullBranchStrategy() string {
	return getConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := GetConfigurationValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return command.New("git", "remote", "get-url", "origin").Output()
}

// GetRemoteUpstreamURL returns the URL of the "upstream" remote.
func GetRemoteUpstreamURL() string {
	return command.New("git", "remote", "get-url", "upstream").Output()
}

// GetURLHostname returns the hostname contained within the given Git URL.
func GetURLHostname(url string) string {
	hostnameRegex, err := regexp.Compile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	exit.IfWrap(err, "Error compiling hostname regular expression")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func GetURLRepositoryName(url string) string {
	hostname := GetURLHostname(url)
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
	return command.New("git", "config", "-l", "--global", "--name").OutputContainsLine(key)
}

// HasParentBranch returns whether or not the given branch has a parent
func HasParentBranch(branchName string) bool {
	return GetParentBranch(branchName) != ""
}

// IsAncestorBranch returns whether the given branch is an ancestor of the other given branch.
func IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := GetAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func HasRemote(name string) bool {
	return util.DoesStringArrayContain(getRemotes(), name)
}

// IsFeatureBranch returns whether the branch with the given name is
// a feature branch.
func IsFeatureBranch(branchName string) bool {
	return !IsMainBranch(branchName) && !IsPerennialBranch(branchName)
}

// IsMainBranch returns whether the branch with the given name
// is the main branch of the repository.
func IsMainBranch(branchName string) bool {
	return branchName == GetMainBranch()
}

// IsPerennialBranch returns whether the branch with the given name is
// a perennial branch.
func IsPerennialBranch(branchName string) bool {
	perennialBranches := GetPerennialBranches()
	return util.DoesStringArrayContain(perennialBranches, branchName)
}

// RemoveAllConfiguration removes all Git Town configuration
func RemoveAllConfiguration() {
	command.New("git", "config", "--remove-section", "git-town").Output()
}

// RemoveOutdatedConfiguration removes outdated Git Town configuration
func RemoveOutdatedConfiguration() {
	for child, parent := range GetParentBranchMap() {
		if !HasBranch(child) || !HasBranch(parent) {
			DeleteParentBranch(child)
		}
	}
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch
func RemoveFromPerennialBranches(branchName string) {
	SetPerennialBranches(util.RemoveStringFromSlice(GetPerennialBranches(), branchName))
}

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func SetMainBranch(branchName string) {
	setConfigurationValue("git-town.main-branch-name", branchName)
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func SetParentBranch(branchName, parentBranchName string) {
	setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// SetPerennialBranches marks the given branches as perennial branches
func SetPerennialBranches(branchNames []string) {
	setConfigurationValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

// SetPullBranchStrategy updates the configured pull branch strategy.
func SetPullBranchStrategy(strategy string) {
	setConfigurationValue("git-town.pull-branch-strategy", strategy)
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

func getConfigurationKeysMatching(toMatch string) (result []string) {
	re, err := regexp.Compile(toMatch)
	exit.IfWrapf(err, "Error compiling configuration regular expression (%s): %v", toMatch, err)
	return configMap.KeysMatching(re)
}

func setConfigurationValue(key, value string) {
	command.New("git", "config", key, value).Run()
	configMap.Set(key, value)
}

func setGlobalConfigurationValue(key, value string) {
	command.New("git", "config", "--global", key, value).Run()
	globalConfigMap.Set(key, value)
	configMap.Reset() // Need to reset config in case it was inheriting
}

func removeConfigurationValue(key string) {
	command.New("git", "config", "--unset", key).Run()
	configMap.Delete(key)
}

// Remotes are cached in order to minimize the number of git commands run
var remotes []string
var remotesInitialized bool

func getRemotes() []string {
	if !remotesInitialized {
		remotes = strings.Split(command.New("git", "remote").Output(), "\n")
		remotesInitialized = true
	}
	return remotes
}

// Init

func init() {
	configMap = NewConfigMap(false)
	globalConfigMap = NewConfigMap(true)
}
