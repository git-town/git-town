package config

import (
	"os"
	"strings"

	"github.com/Originate/git-town/lib/util"
)

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

func DeleteAncestorBranches(branchName string) {
	removeConfigurationValue("git-town-branch." + branchName + ".ancestors")
}

func DeleteParentBranch(branchName string) {
	removeConfigurationValue("git-town-branch." + branchName + ".parent")
}

func GetAncestorBranches(branchName string) []string {
	value := getConfigurationValue("git-town-branch." + branchName + ".ancestors")
	if value == "" {
		return []string{}
	}
	return strings.Split(value, " ")
}

func GetChildBranches(branchName string) (result []string) {
	configs := util.GetCommandOutput("git", "config", "--get-regexp", "^git-town-branch\\..*\\.parent$")
	for _, config := range strings.Split(configs, "\n") {
		splitConfig := strings.Split(config, " ")
		key := splitConfig[0]
		value := splitConfig[1]
		if value == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append([]string{child}, result...)
		}
	}
	return
}

func GetMainBranch() string {
	return getConfigurationValue("git-town.main-branch-name")
}

func GetParentBranch(branchName string) string {
	return getConfigurationValue("git-town-branch." + branchName + ".parent")
}

func GetPerennialBranches() []string {
	return strings.Split(getConfigurationValue("git-town.perennial-branch-names"), " ")
}

func GetPullBranchStrategy() string {
	return getConfigurationValueWithDefault("git-town.pull-branch-strategy", "rebase")
}

func GetRemoteOriginUrl() string {
	if os.Getenv("GIT_TOWN_ENV") == "test" {
		mockRemoteUrl := getConfigurationValue("git-town.testing.remote-url")
		if mockRemoteUrl != "" {
			return mockRemoteUrl
		}
	}
	return util.GetCommandOutput("git", "remote", "get-url", "origin")
}

func GetRemoteUpstreamUrl() string {
	return util.GetCommandOutput("git", "remote", "get-url", "upstream")
}

func IsAncestorBranch(branchName, ancestorBranchName string) bool {
	ancestorBranches := CompileAncestorBranches(branchName)
	return util.DoesStringArrayContain(ancestorBranches, ancestorBranchName)
}

func HasCompiledAncestorBranches(branchName string) bool {
	return len(GetAncestorBranches(branchName)) > 0
}

func HasRemote(name string) bool {
	return util.DoesCommandOuputContainLine([]string{"git", "remote"}, name)
}

func IsFeatureBranch(branchName string) bool {
	return !IsMainBranch(branchName) && !IsPerennialBranch(branchName)
}

func IsMainBranch(branchName string) bool {
	return branchName == GetMainBranch()
}

func IsPerennialBranch(branchName string) bool {
	perennialBranches := GetPerennialBranches()
	return util.DoesStringArrayContain(perennialBranches, branchName)
}

func SetAncestorBranches(branchName string, ancestorBranches []string) {
	setConfigurationValue("git-town-branch."+branchName+".ancestors", strings.Join(ancestorBranches, " "))
}

func SetParentBranch(branchName, parentBranchName string) {
	setConfigurationValue("git-town-branch."+branchName+".parent", parentBranchName)
}

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

func hasConfigurationValue(key string) bool {
	return util.DoesCommandOuputContainLine([]string{"git", "config", "-l", "--local", "--name"}, key)
}

func setConfigurationValue(key, value string) {
	util.GetCommandOutput("git", "config", key, value)
}

func removeConfigurationValue(key string) {
	util.GetCommandOutput("git", "config", "--unset", key)
}
