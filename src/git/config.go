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

// AddToPerennialBranches adds the given branch as a perennial branch
func AddToPerennialBranches(branchName string) {
	SetPerennialBranches(append(GetPerennialBranches(), branchName))
}

// CompileAncestorBranches calculates and returns the list of ancestor branches
// of the given branch based off the "git-town-branch.XXX.parent" configuration values.
func CompileAncestorBranches(branchName string) (result []string) {
	current := branchName
	for {
		if IsMainBranch(current) || IsPerennialBranch(current) {
			return
		}
		parent := GetParentBranch(current)
		if parent == "" {
			return
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// DeleteAllAncestorBranches removes all Git Town ancestor entries
// for all branches from the configuration.
func DeleteAllAncestorBranches() {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.ancestors$") {
		removeConfigurationValue(key)
	}
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
// beginning but not including the parennial branch from which this hierarchy was cut.
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func GetAncestorBranches(branchName string) []string {
	value := getLocalConfigurationValue("git-town-branch." + branchName + ".ancestors")
	if value == "" {
		return []string{}
	}
	return strings.Split(value, " ")
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func GetChildBranches(branchName string) (result []string) {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := getLocalConfigurationValue(key)
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
	return command.New("git", "config", key).Output()
}

// GetGlobalConfigurationValue returns the global git configuration value for the given key
func GetGlobalConfigurationValue(key string) (result string) {
	if hasConfigurationValue("global", key) {
		result = command.New("git", "config", "--global", key).Output()
	}
	return
}

// GetMainBranch returns the name of the main branch.
func GetMainBranch() string {
	return getLocalConfigurationValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func GetParentBranch(branchName string) string {
	return getLocalConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func GetPerennialBranches() []string {
	result := getLocalConfigurationValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
func GetPullBranchStrategy() string {
	return getLocalConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := getLocalConfigurationValue("git-town.testing.remote-url")
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

// IsAncestorBranch returns whether the given branch is an ancestor of the other given branch.
func IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := CompileAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

// HasCompiledAncestorBranches returns whether the Git Town configuration
// contains a cached ancestor list for the branch with the given name.
func HasCompiledAncestorBranches(branchName string) bool {
	return len(GetAncestorBranches(branchName)) > 0
}

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func HasRemote(name string) bool {
	return command.New("git", "remote").OutputContainsLine(name)
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

// RemoveFromPerennialBranches removes the given branch as a perennial branch
func RemoveFromPerennialBranches(branchName string) {
	SetPerennialBranches(util.RemoveStringFromSlice(GetPerennialBranches(), branchName))
}

// SetAncestorBranches stores the given list of branches as ancestors
// for the given branch in the Git Town configuration.
func SetAncestorBranches(branchName string, ancestorBranches []string) {
	setConfigurationValue("git-town-branch."+branchName+".ancestors", strings.Join(ancestorBranches, " "))
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

// UpdateOffline updates whether Git Town is in offline mode
func UpdateOffline(value bool) {
	setGlobalConfigurationValue("git-town.offline", strconv.FormatBool(value))
}

// UpdateShouldNewBranchPush updates whether the current repository is configured to push
// freshly created branches up to the origin remote.
func UpdateShouldNewBranchPush(value bool) {
	setConfigurationValue("git-town.new-branch-push-flag", strconv.FormatBool(value))
}

// Helpers

func getConfigurationValueWithDefault(key, defaultValue string) string {
	value := GetConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getLocalConfigurationValue returns the given configuration value
// only from the local Git configuration
func getLocalConfigurationValue(key string) (result string) {
	if hasConfigurationValue("local", key) {
		result = command.New("git", "config", "--local", key).Output()
	}
	return
}

func getLocalConfigurationValueWithDefault(key, defaultValue string) string {
	value := getLocalConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getConfigurationKeysMatching(toMatch string) (result []string) {
	configRegexp, err := regexp.Compile(toMatch)
	exit.IfWrapf(err, "Error compiling configuration regular expression (%s): %v", toMatch, err)
	lines := command.New("git", "config", "-l", "--local", "--name").Output()
	for _, line := range strings.Split(lines, "\n") {
		if configRegexp.MatchString(line) {
			result = append(result, line)
		}
	}
	return
}

func hasConfigurationValue(location, key string) bool {
	return command.New("git", "config", "-l", "--"+location, "--name").OutputContainsLine(key)
}

func setConfigurationValue(key, value string) {
	command.New("git", "config", key, value).Run()
}

func setGlobalConfigurationValue(key, value string) {
	command.New("git", "config", "--global", key, value).Run()
}

func removeConfigurationValue(key string) {
	command.New("git", "config", "--unset", key).Run()
}
