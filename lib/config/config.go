package config

import (
  "os"
  "sort"
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
  return sort.SearchStrings(perennialBranches, branchName) < len(perennialBranches)
}


func HasRemoteOrigin() bool {
  return GetRemoteOriginUrl() != ""
}


func HasRemoteUpstream() bool {
  return GetRemoteUpstreamUrl() != ""
}


func SetParentBranch(branchName, parentBranchName string) {
  storeConfigurationValue(branchName + ".parent", parentBranchName)
}

func ShouldHackPush() bool {
  return getConfigurationValueWithDefault("hack-push-flag", "true") == "true"
}


// Helpers


func getConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", "git-town." + key})
}

func getConfigurationValueWithDefault(key, defaultValue string) string {
  value := getConfigurationValue(key)
  if value == "" {
    return defaultValue
  }
  return value
}

func storeConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", "git-town." + key, value})
}
