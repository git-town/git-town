package config

import (
  "os"
  "strings"

  "github.com/Originate/gt/lib/util"
)


func GetMainBranch() string {
  return getConfigurationValue("main-branch-name")
}


func GetParentBranch(branchName string) string {
  return getConfigurationValue(branchName + ".parent")
}


func GetPerennialBranches() []string {
  return strings.Split(getConfigurationValue("perennial-branch-names"), " ")
}


func GetPullBranchStrategy() string {
  return getConfigurationValueWithDefault("pull-branch-strategy", "rebase")
}


func GetRemoteOriginUrl() string {
  if os.Getenv("GIT_TOWN_ENV") == "test" {
    mockRemoteUrl := getConfigurationValue("testing.remote-url")
    if mockRemoteUrl != "" {
      return mockRemoteUrl
    }
  }
  return util.GetCommandOutput([]string{"git", "remote", "get-url", "origin"})
}


func GetRemoteUpstreamUrl() string {
  return util.GetCommandOutput([]string{"git", "remote", "get-url", "upstream"})
}


func IsFeatureBranch(branchName string) bool {
  return branchName != GetMainBranch() && !IsPerennialBranch(branchName)
}


func IsPerennialBranch(branchName string) bool {
  perennialBranches := GetPerennialBranches()
  return util.DoesStringArrayContain(perennialBranches, branchName)
}


func HasRemoteOrigin() bool {
  return hasRemote("origin")
}


func HasRemoteUpstream() bool {
  return hasRemote("upstream")
}


func SetParentBranch(branchName, parentBranchName string) {
  storeConfigurationValue(branchName + ".parent", parentBranchName)
}


func ShouldHackPush() bool {
  return getConfigurationValueWithDefault("hack-push-flag", "true") == "true"
}


// Helpers


func getConfigurationValue(key string) string {
  namespacedKey := "git-town." + key
  value := ""
  if hasConfigurationValue(namespacedKey) {
    value = util.GetCommandOutput([]string{"git", "config", namespacedKey})
  }
  return value
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


func hasRemote(name string) bool {
  return util.DoesCommandOuputContainLine([]string{"git", "remote"}, name)
}


func storeConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", "git-town." + key, value})
}
