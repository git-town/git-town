package config

import (
  "os"
  "strings"

  "github.com/Originate/git-town/lib/util"
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
  return util.GetCommandOutput("git", "remote", "get-url", "origin")
}


func GetRemoteUpstreamUrl() string {
  return util.GetCommandOutput("git", "remote", "get-url", "upstream")
}


func IsFeatureBranch(branchName string) bool {
  return branchName != GetMainBranch() && !IsPerennialBranch(branchName)
}


func IsPerennialBranch(branchName string) bool {
  perennialBranches := GetPerennialBranches()
  return util.DoesStringArrayContain(perennialBranches, branchName)
}

func HasRemote(name string) bool {
  return util.DoesCommandOuputContainLine([]string{"git", "remote"}, name)
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
    value = util.GetCommandOutput("git", "config", namespacedKey)
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


func storeConfigurationValue(key, value string) {
  util.GetCommandOutput("git", "config", "git-town." + key, value)
}
