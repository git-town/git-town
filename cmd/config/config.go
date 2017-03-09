package config

import (
  "os"
  "sort"
  "strings"

  "github.com/Originate/gt/cmd/util"
)


func CompileAncestorBranches(branchName string) (result []string) {
  current := branchName
  for {
    parent := GetParentBranch(current)
    result = append(result, parent)
    if parent == GetMainBranch() {
      break
    }
    current = parent
  }
  return
}


func GetAncestorBranches(branchName string) []string {
  value := getBranchConfigurationValue(branchName + ".ancestors")
  if value == "" {
    return []string{}
  }
  return strings.Split(value, " ")
}


func GetMainBranch() string {
  return getGlobalConfigurationValue("main-branch-name")
}


func GetParentBranch(branchName string) string {
  return getBranchConfigurationValue(branchName + ".parent")
}


func GetPerennialBranches() []string {
  return strings.Split(getGlobalConfigurationValue("perennial-branch-names"), " ")
}


func GetPullBranchStrategy() string {
  return getGlobalConfigurationValueWithDefault("pull-branch-strategy", "rebase")
}


func GetRemoteOriginUrl() string {
  if os.Getenv("GIT_TOWN_ENV") == "test" {
    mockRemoteUrl := getGlobalConfigurationValue("testing.remote-url")
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

func KnowsAllAncestorBranches(branchName string) bool {
  return branchName == GetMainBranch() ||
    IsPerennialBranch(branchName) ||
    len(GetAncestorBranches(branchName)) > 0
}

func HasRemoteOrigin() bool {
  return GetRemoteOriginUrl() != ""
}


func HasRemoteUpstream() bool {
  return GetRemoteUpstreamUrl() != ""
}


func SetAncestorBranches(branchName string, ancestorBranches []string) {
  setBranchConfigurationValue(branchName + ".ancestors", strings.Join(ancestorBranches, " "))
}


func SetParentBranch(branchName, parentBranchName string) {
  setBranchConfigurationValue(branchName + ".parent", parentBranchName)
}


func ShouldHackPush() bool {
  return getGlobalConfigurationValueWithDefault("hack-push-flag", "true") == "true"
}


// Helpers

func getConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", key})
}

func getBranchConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", "git-town-branch." + key})
}

func getGlobalConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", "git-town." + key})
}

func getGlobalConfigurationValueWithDefault(key, defaultValue string) string {
  value := getGlobalConfigurationValue(key)
  if value == "" {
    return defaultValue
  }
  return value
}

func setConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", key, value})
}

func setBranchConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", "git-town-branch." + key, value})
}

func setGlobalConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", "git-town." + key, value})
}
