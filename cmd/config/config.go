package config

import (
  "os"
  "strings"

  "github.com/Originate/gt/cmd/util"
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
  pullBranchStrategy := getConfigurationValueWithDefault("pull-branch-strategy", "rebase")
  return pullBranchStrategy
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
  return util.ContainsString(perennialBranches, branchName)
}


func HasRemoteOrigin() bool {
  return GetRemoteOriginUrl() != ""
}


func HasRemoteUpstream() bool {
  return GetRemoteUpstreamUrl() != ""
}


func ShouldHackPush() bool {
  hackPushFlag := getConfigurationValueWithDefault("hack-push-flag", "true")
  return hackPushFlag == "true"
}


func StoreParentBranch(branchName, parentBranchName string) {
  storeConfigurationValue(branchName + ".parent", parentBranchName)
}


// Helpers


func getConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", "git-town." + key})
}

func getConfigurationValueWithDefault(key, defaultValue string) string {
  value := util.GetCommandOutput([]string{"git", "config", "git-town." + key})
  if value == "" {
    return defaultValue
  }
  return value
}

func storeConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", "git-town." + key, value})
}
