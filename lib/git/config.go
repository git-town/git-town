/*
This file contains functionality around storing configuration settings
inside Git's metadata storage for the repository.
*/

package git

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Originate/git-town/lib/util"
)

// AddPerennialBranch adds the given branch as a perennial branch
func AddPerennialBranch(branchName string) {
	setPerennialBranches(append(GetPerennialBranches(), branchName))
}

// CompileAncestorBranches re-calculates and returns the list of ancestor branches
// of the given branch
// based off the "git-town-branch.XXX.ancestors" configuration values.
// The result starts with but does not include the perennial branch
// from which this branch hierarchy was cut initially.
func CompileAncestorBranches(branchName string) (result []string) {
	current := branchName
	for {
		parent := GetParentBranch(current)
		result = append([]string{parent}, result...)
		if IsMainBranch(parent) || IsPerennialBranch(parent) {
			return
		}
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
	value := getConfigurationValue("git-town-branch." + branchName + ".ancestors")
	if value == "" {
		return []string{}
	}
	return strings.Split(value, " ")
}

// GetChildBranches returns the names of all branches for which the given branch
// is a parent.
func GetChildBranches(branchName string) (result []string) {
	for _, key := range getConfigurationKeysMatching("^git-town-branch\\..*\\.parent$") {
		parent := getConfigurationValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append([]string{child}, result...)
		}
	}
	return
}

// GetMainBranch returns the name of the main branch.
func GetMainBranch() string {
	return getConfigurationValue("git-town.main-branch-name")
}

// GetParentBranch returns the name of the parent branch of the given branch.
func GetParentBranch(branchName string) string {
	return getConfigurationValue("git-town-branch." + branchName + ".parent")
}

// GetPerennialBranches returns all branches that are marked as perennial.
func GetPerennialBranches() []string {
	return strings.Split(getConfigurationValue("git-town.perennial-branch-names"), " ")
}

// GetPullBranchStrategy returns the currently configured pull branch strategy.
// See https://github.com/Originate/git-town/blob/master/documentation/commands/git-town.md
func GetPullBranchStrategy() string {
	return getConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

// GetRemoteOriginURL returns the URL for the "origin" remote.
// In tests this value can be stubbed.
func GetRemoteOriginURL() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteURL := getConfigurationValue("git-town.testing.remote-url")
		if mockRemoteURL != "" {
			return mockRemoteURL
		}
	}
	return util.GetCommandOutput("git", "remote", "get-url", "origin")
}

// GetRemoteUpstreamURL returns the URL of the "upstream" remote.
func GetRemoteUpstreamURL() string {
	return util.GetCommandOutput("git", "remote", "get-url", "upstream")
}

// GetURLHostname returns the hostname contained within the given Git URL.
func GetURLHostname(url string) string {
	hostnameRegex, err := regexp.Compile("(^[^:]*://([^@]*@)?|git@)([^/:]+).*")
	if err != nil {
		log.Fatal("Error compiling hostname regular expression: ", err)
	}
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
	if err != nil {
		log.Fatal("Error compiling repository name regular expression: ", err)
	}
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
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
	return util.DoesCommandOuputContainLine([]string{"git", "remote"}, name)
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

// RemovePerennialBranch removes the given branch as a perennial branch
func RemovePerennialBranch(branchName string) {
	setPerennialBranches(util.RemoveStringFromSlice(GetPerennialBranches(), branchName))
}

// SetAncestorBranches stores the given list of branches as ancestors
// for the given branch in the Git Town configuration.
func SetAncestorBranches(branchName string, ancestorBranches []string) {
	setConfigurationValue("git-town-branch."+branchName+".ancestors", strings.Join(ancestorBranches, " "))
}

// SetParentBranch marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func SetParentBranch(branchName, parentBranchName string) {
	setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

// ShouldHackPush returns whether the current repository is configured to push
// freshly created branches up to the origin remote.
func ShouldHackPush() bool {
	return getConfigurationValueWithDefault("git-town.hack-push-flag", "true") == "true"
}

// Helpers

func getConfigurationValue(key string) (result string) {
	if hasConfigurationValue(key) {
		result = util.GetCommandOutput("git", "config", key)
	}
	return
}

func getConfigurationValueWithDefault(key, defaultValue string) string {
	value := getConfigurationValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getConfigurationKeysMatching(toMatch string) (result []string) {
	configRegexp, err := regexp.Compile(toMatch)
	if err != nil {
		log.Fatalf("Error compiling configuration regular expression (%s): %v", toMatch, err)
	}
	lines := util.GetCommandOutput("git", "config", "-l", "--local", "--name")
	for _, line := range strings.Split(lines, "\n") {
		if configRegexp.MatchString(line) {
			result = append(result, line)
		}
	}
	return
}

func hasConfigurationValue(key string) bool {
	return util.DoesCommandOuputContainLine([]string{"git", "config", "-l", "--local", "--name"}, key)
}

func setConfigurationValue(key, value string) {
	util.GetCommandOutput("git", "config", key, value)
}

func setPerennialBranches(branchNames []string) {
	setConfigurationValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
}

func removeConfigurationValue(key string) {
	util.GetCommandOutput("git", "config", "--unset", key)
}
