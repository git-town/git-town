package config

import (
  "fmt"
  "os"
  "strings"

  "github.com/Originate/gt/cmd/util"
)


func IsFeatureBranch(branchName string) bool {
  return branchName != GetMainBranch() && !IsPernnialBranch(branchName)
}

func IsPernnialBranch(branchName string) bool {
  perennialBranches := GetPerennialBranches()
  return util.Contains(perennialBranches, branchName)
}

func GetParentBranch(branchName string) string {
  return getConfigurationValue(fmt.Sprintf("%s.parent", branchName))
}

func GetMainBranch() string {
  return getConfigurationValue("main-branch-name")
}

func GetPerennialBranches() []string {
  return strings.Split(getConfigurationValue("perennial-branch-names"), " ")
}

func GetPullBranchStrategy() string {
  pullBranchStrategy := getConfigurationValue("pull-branch-strategy")
  if pullBranchStrategy == "" {
    pullBranchStrategy = "rebase"
  }
  return pullBranchStrategy
}

func GetRemoteUrl() string {
  if os.Getenv("GIT_TOWN_ENV") == "test" {
    mockRemoteUrl := getConfigurationValue("testing.remote-url")
    if mockRemoteUrl != "" {
      return mockRemoteUrl
    }
  }
  return util.GetCommandOutput([]string{"git", "remote", "get-url", "origin"})
}

func HasRemote() bool {
  return GetRemoteUrl() != ""
}

func ShouldHackPush() bool {
  hackPushFlag := getConfigurationValue("hack-push-flag")
  if hackPushFlag == "" {
    hackPushFlag = "true"
  }
  return hackPushFlag == "true"
}

func StoreParentBranch(branchName, parentBranchName string) {
  storeConfigurationValue(fmt.Sprintf("%s.parent", branchName), parentBranchName)
}

func getConfigurationValue(key string) string {
  return util.GetCommandOutput([]string{"git", "config", fmt.Sprintf("git-town.%s", key)})
}

func storeConfigurationValue(key, value string) {
  util.GetCommandOutput([]string{"git", "config", fmt.Sprintf("git-town.%s", key), value})
}
